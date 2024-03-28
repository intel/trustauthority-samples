/*
 * Copyright (c) 2024 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/intel/trustauthority-samples/tdxexample/model"
	"github.com/intel/trustauthority-samples/tdxexample/service"
	httpTransport "github.com/intel/trustauthority-samples/tdxexample/transport/http"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	ServiceDir         = "trustauthority-demo/"
	HomeDir            = "/opt/" + ServiceDir
	DefaultTLSCertPath = HomeDir + "tls.crt"
	DefaultTLSKeyPath  = HomeDir + "tls.key"
	ValidityDays       = 365
	DefaultKeyLength   = 3072
)

func main() {
	conf, err := configure()
	if err != nil {
		panic(err)
	}

	err = conf.Validate()
	if err != nil {
		panic(err)
	}

	// Initialize http client
	tlsConfig := &tls.Config{}
	if conf.SkipTLSVerification {
		tlsConfig.InsecureSkipVerify = true
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Initialize user data
	privKey, pubBytes, err := generateKeyPair()
	if err != nil {
		panic(err)
	}
	userData := base64.StdEncoding.EncodeToString(pubBytes)

	// Initialize Model Executor
	modelExecutor := model.NewModelExecutor("/etc/model.enc", privKey)

	// Initialize the Service
	svc, err := service.NewService(userData, httpClient, modelExecutor)
	if err != nil {
		panic(err)
	}

	// Associate the service to rest endpoints/http
	httpHandlers, err := httpTransport.NewHTTPHandler(svc)
	if err != nil {
		panic(err)
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", conf.Port),
		Handler:           httpHandlers,
		ReadHeaderTimeout: time.Duration(conf.HTTPReadHdrTimeout) * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
			CipherSuites: []uint16{tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_AES_128_GCM_SHA256,
				// TLS_AES_128_CCM_SHA256 is not supported by go crypto/tls package
				tls.TLS_CHACHA20_POLY1305_SHA256},
		},
	}

	// TLS certificate is passed
	if _, err := os.Stat(DefaultTLSCertPath); os.IsNotExist(err) {
		// TLS certificate and key does not exist, so creating the cert and key
		err = generateTLSKeyandCert(DefaultTLSCertPath, DefaultTLSKeyPath, conf.SanList)
		if err != nil {
			panic(err)
		}
	}
	log.Debugf("Starting HTTPS server with TLS cert: %s", DefaultTLSCertPath)

	if err := httpServer.ListenAndServeTLS(DefaultTLSCertPath, DefaultTLSKeyPath); err != nil {
		panic(err)
	}
}

func generateKeyPair() (*rsa.PrivateKey, []byte, error) {
	keyPair, err := rsa.GenerateKey(rand.Reader, DefaultKeyLength)
	defer ZeroizeRSAPrivateKey(keyPair)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error while generating RSA key pair")
	}

	// Public key format : <exponent:E_SIZE_IN_BYTES><modulus:N_SIZE_IN_BYTES>
	pub := keyPair.PublicKey
	pubBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(pubBytes, uint32(pub.E))
	pubBytes = append(pubBytes, pub.N.Bytes()...)
	return keyPair, pubBytes, nil
}

func generateTLSKeyandCert(TLSCertPath, TLSKeyPath, TlsSanList string) error {
	key, err := rsa.GenerateKey(rand.Reader, DefaultKeyLength)
	defer ZeroizeRSAPrivateKey(key)
	if err != nil {
		return errors.Wrap(err, "error while generating RSA key pair")
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return errors.Wrap(err, "Failed to create serial number")
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},

		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, ValidityDays),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// add the san list for tls certificate
	hosts := strings.Split(TlsSanList, ",")
	for _, h := range hosts {
		h = strings.TrimSpace(h)
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	// store key and cert to file
	selfSignCert, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	defer ZeroizeByteArray(selfSignCert)
	if err != nil {
		return errors.Wrap(err, "Failed to create certificate")
	}
	tlsCertPath := filepath.Clean(TLSCertPath)
	certOut, err := os.OpenFile(tlsCertPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|syscall.O_NOFOLLOW, 0600)
	if err != nil {
		return fmt.Errorf("could not open %s file for writing: %v", tlsCertPath, err)
	}
	defer func() {
		derr := certOut.Close()
		if derr != nil {
			log.WithError(derr).Error("Error closing Cert file")
		}
	}()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: selfSignCert}); err != nil {
		return fmt.Errorf("could not pem encode cert: %v", err)
	}

	keyDer, err := x509.MarshalPKCS8PrivateKey(key)
	defer ZeroizeByteArray(keyDer)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal private key")
	}
	tlsKeyPath := filepath.Clean(TLSKeyPath)
	keyOut, err := os.OpenFile(tlsKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|syscall.O_NOFOLLOW, 0600) // open file with restricted permissions
	if err != nil {
		return fmt.Errorf("could not open %s file for writing: %v", tlsKeyPath, err)
	}
	defer func() {
		derr := keyOut.Close()
		if derr != nil {
			log.WithError(derr).Error("Error closing Key file")
		}
	}()
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: keyDer}); err != nil {
		return fmt.Errorf("could not pem encode the private key: %v", err)
	}

	return nil
}
