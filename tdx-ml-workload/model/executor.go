/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package model

// #cgo CFLAGS: -fno-strict-overflow -fno-delete-null-pointer-checks -fwrapv -fstack-protector-strong
// #cgo CXXFLAGS: -I/usr/lib/
// #cgo LDFLAGS: -L/usr/lib/ -lstdc++ -lcrypto
// #include "model.h"
import "C"

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"os"
	"unsafe"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	privateKeyLocation = "envelope.key"
)

type ModelExecutor struct {
	modelPath string
}

func NewModelExecutor(mpath string) *ModelExecutor {
	return &ModelExecutor{
		modelPath: mpath,
	}
}

func (m *ModelExecutor) ResetModel() error {
	log.Debug("Resetting Model. ")

	status := C.model_reset()
	if status != 0 {
		log.Errorf("Resetting AI model failed! Error code: 0x%04x", status)
		return errors.New("Resetting AI model failed!")
	}

	return nil
}

func (m *ModelExecutor) ExecuteModel(pregnancies float32,
	glucose float32, bloodpressure float32, skinthickness float32, insulin float32, bmi float32,
	dbf float32, age float32) (int, error) {

	log.Debug("Executing Model. ")

	var inferedValue C.int
	status := C.model_predict(C.double(pregnancies),
		C.double(glucose),
		C.double(bloodpressure),
		C.double(skinthickness),
		C.double(insulin),
		C.double(bmi),
		C.double(dbf),
		C.double(age),
		(*C.int)(&inferedValue))
	if status != 0 {
		log.Errorf("AI Model Inferencing failed! Error code: 0x%04x", status)
		return -1, errors.New("AI Model Inferencing failed!")
	}

	log.Debugf("Golang Diabetes Prediction : %d", inferedValue)
	return int(inferedValue), nil
}

func (m *ModelExecutor) DecryptModel(swk, dek []byte) error {
	log.Debug("Decrypting Model. ")
	// read ai model from a file
	model, err := os.ReadFile(m.modelPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("AI model file does not exist. Please check filepath again.")
		} else {
			return errors.Wrapf(err, "Unable to read ai model from file : %s", m.modelPath)
		}
	}

	if len(model) == 0 {
		log.Error("Size of ai model can't be zero!")
		return errors.New("Size of ai model can't be zero!")
	}

	key, err := UnwrapKey(swk, privateKeyLocation)
	if err != nil {
		return errors.Wrap(err, "Error while unwrapping the swk")
	}

	//Decrypt the model inside the TD
	status := C.model_decrypt((*C.uint8_t)(unsafe.Pointer(&model[0])),
		C.ulong(len(model)),
		(*C.uint8_t)(unsafe.Pointer(&dek[12])),
		C.ulong(len(dek[12:])),
		(*C.uint8_t)(unsafe.Pointer(&key[0])))
	if status != 0 {
		log.Errorf("Decryption of model failed! Error code: 0x%04x", status)
		return errors.New("Decryption of model failed!")
	}

	return nil
}

func UnwrapKey(wrappedKey []byte, privateKeyLocation string) ([]byte, error) {

	privateKey, err := os.ReadFile(privateKeyLocation)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading private envelope key file")
	}

	privateKeyBlock, _ := pem.Decode(privateKey)
	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Error decoding private envelope key")
	}

	decryptedKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pri, wrappedKey, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error while decrypting the swk")
	}

	log.Info("Successfully unwrapped swk")
	return decryptedKey, nil
}
