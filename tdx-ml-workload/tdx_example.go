/*
 * Copyright (c) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	"github.com/intel/trustauthority-samples/tdxexample/model"
	"github.com/intel/trustauthority-samples/tdxexample/service"
	httpTransport "github.com/intel/trustauthority-samples/tdxexample/transport/http"
	"github.com/pkg/errors"
)

const (
	envelopePrivateKeyFile = "envelope.key"
	pemBlockTypePrivateKey = "PRIVATE KEY"
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

	// Initialize Model Executor
	modelExecutor := model.NewModelExecutor("/etc/model.enc")

	// Initialize http client
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Initialize user data
	pubBytes, err := loadPublicKey()
	if err != nil {
		panic(err)
	}
	userData := base64.StdEncoding.EncodeToString(pubBytes)

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
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: httpHandlers,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

func loadPublicKey() ([]byte, error) {
	keyPair, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		return nil, errors.Wrap(err, "error while generating RSA key pair")
	}

	// save private key
	privateKey := &pem.Block{
		Type:  pemBlockTypePrivateKey,
		Bytes: x509.MarshalPKCS1PrivateKey(keyPair),
	}

	privateKeyFile, err := os.OpenFile(envelopePrivateKeyFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return nil, errors.Wrap(err, "I/O error while saving private key")
	}
	defer func() {
		derr := privateKeyFile.Close()
		if derr != nil {
			fmt.Printf("Error while closing file" + derr.Error())
		}
	}()

	err = pem.Encode(privateKeyFile, privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "error while encoding the private key")
	}

	// Public key format : <exponent:E_SIZE_IN_BYTES><modulus:N_SIZE_IN_BYTES>
	pub := keyPair.PublicKey
	pubBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(pubBytes, uint32(pub.E))
	pubBytes = append(pubBytes, pub.N.Bytes()...)
	return pubBytes, nil
}
