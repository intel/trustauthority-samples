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
	"os"
	"path/filepath"
	"unsafe"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ModelExecutor struct {
	modelPath string
	privKey   *rsa.PrivateKey
}

func NewModelExecutor(mpath string, pkey *rsa.PrivateKey) *ModelExecutor {
	return &ModelExecutor{
		modelPath: mpath,
		privKey:   pkey,
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
	modelPath := filepath.Clean(m.modelPath)
	model, err := os.ReadFile(modelPath)
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

	key, err := UnwrapKey(swk, m.privKey)
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

func UnwrapKey(wrappedKey []byte, pri *rsa.PrivateKey) ([]byte, error) {

	decryptedKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pri, wrappedKey, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error while decrypting the swk")
	}

	log.Info("Successfully unwrapped swk")
	return decryptedKey, nil
}
