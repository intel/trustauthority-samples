/*
 *   Copyright (c) 2024 Intel Corporation
 *   All rights reserved.
 *   SPDX-License-Identifier: BSD-3-Clause
 */
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// zeroizeByteArray overwrites a byte array's data with zeros
func zeroizeByteArray(bytes []byte) {
	for i, _ := range bytes {
		bytes[i] = 0
	}
}

// zeroizeBigInt replaces the big integer's byte array with
// zeroes.  This function will panic if the bigInt parameter is nil.
func zeroizeBigInt(bigInt *big.Int) {
	if bigInt == nil {
		panic("The bigInt parameter cannot be nil")
	}

	bytes := make([]byte, len(bigInt.Bytes()))
	bigInt.SetBytes(bytes)
}

// zeroizeRSAPrivateKey clears the private key's "D" and
// "Primes" (big int) values.  This function will panic if the privateKey
// parameter is nil.
func zeroizeRSAPrivateKey(privateKey *rsa.PrivateKey) {
	if privateKey == nil {
		panic("The private key parameter cannot be nil")
	}

	zeroizeBigInt(privateKey.D)
	for _, bigInt := range privateKey.Primes {
		zeroizeBigInt(bigInt)
	}
}

func Encrypt(modelPath string, privateKeyLocation string, encryptedFileLocation string, wrappedKey []byte) error {

	modelPath = filepath.Clean(modelPath)
	model, err := os.ReadFile(modelPath)
	if err != nil {
		return errors.Wrap(err, "Error reading the data file")
	}

	key, err := UnwrapKey(wrappedKey, privateKeyLocation)
	if err != nil {
		return errors.Wrap(err, "Error while unwrapping the key")
	}
	defer zeroizeByteArray(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.Wrap(err, "Error initializing cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return errors.Wrap(err, "Error creating a cipher block")
	}

	iv := make([]byte, gcm.NonceSize())
	// reading random value into the byte array
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return errors.Wrap(err, "Error creating random IV value")
	}

	encryptedData := gcm.Seal(iv, iv, model, nil)
	err = os.WriteFile(encryptedFileLocation, encryptedData, 0600)
	if err != nil {
		return errors.Wrap(err, "Error during writing the encrypted data to file")
	}

	logrus.Info("Successfully encrypted data")
	return nil
}

func UnwrapKey(wrappedKey []byte, privateKeyLocation string) ([]byte, error) {

	privateKeyLocation = filepath.Clean(privateKeyLocation)
	privateKey, err := os.ReadFile(privateKeyLocation)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading private key file")
	}
	defer zeroizeByteArray(privateKey)

	privateKeyBlock, _ := pem.Decode(privateKey)
	if privateKeyBlock == nil {
		return nil, errors.New("private key not found")
	}
	defer zeroizeByteArray(privateKeyBlock.Bytes)

	pri, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Error decoding private key")
	}
	defer zeroizeRSAPrivateKey(pri.(*rsa.PrivateKey))

	decryptedKey, err := rsa.DecryptOAEP(sha512.New384(), rand.Reader, pri.(*rsa.PrivateKey), wrappedKey, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error while decrypting the key")
	}

	logrus.Info("Successfully unwrapped key")
	return decryptedKey, nil
}

func main() {

	if len(os.Args) < 4 {
		err := errors.New("Invalid usage of encrypt")
		fmt.Println("./encrypt data-file private-key-file wrapped-key-file: ", err)
		os.Exit(1)
	}

	wrappedKeyLocation := filepath.Clean(os.Args[3])
	wrappedKey, err := os.ReadFile(wrappedKeyLocation)
	if err != nil {
		fmt.Println("Error reading wrapped key file: ", err)
		os.Exit(1)
	}

	key, err := base64.StdEncoding.DecodeString(string(wrappedKey))
	if err != nil {
		fmt.Println("Error decoding the data encryption key: ", err)
		os.Exit(1)
	}

	// encrypt the model with key retrieved from KBS
	err = Encrypt(os.Args[1], os.Args[2], "model.enc", key)
	if err != nil {
		fmt.Println("Data encryption failed: ", err)
		os.Exit(1)
	}
}
