/*
 * Copyright (C) 2024 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package model

// #cgo CFLAGS: -fno-strict-overflow -fno-delete-null-pointer-checks -fwrapv -fstack-protector-strong
// #include "model.h"
import "C"

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"os"
	"path/filepath"

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
		log.Errorf("Resetting ML model failed! Error code: 0x%04x", status)
		return errors.New("Resetting ML model failed!")
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
		log.Errorf("ML Model Inferencing failed! Error code: 0x%04x", status)
		return -1, errors.New("ML Model Inferencing failed!")
	}

	log.Debugf("Golang Diabetes Prediction : %d", inferedValue)
	return int(inferedValue), nil
}

func (m *ModelExecutor) DecryptModel(wrappedSwk, wrappedDek []byte) error {
	log.Debug("Decrypting Model. ")
	// read ml model from a file
	modelPath := filepath.Clean(m.modelPath)
	cipherModel, err := os.ReadFile(modelPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("ML model file does not exist. Please check filepath again.")
		} else {
			return errors.Wrapf(err, "Unable to read ml model from file : %s", m.modelPath)
		}
	}

	if len(cipherModel) == 0 {
		return errors.New("Size of ml model can't be zero!")
	}

	swk, err := UnwrapKey(wrappedSwk, m.privKey)
	if err != nil {
		return errors.Wrap(err, "Error while unwrapping the swk")
	}
	log.Debug("Successfully unwrapped swk")

	dek, err := Decrypt(swk, wrappedDek[12:])
	if err != nil {
		return errors.Wrap(err, "Error while decrypting the dek")
	}
	log.Debug("Successfully decrypted dek")

	model, err := Decrypt(dek, cipherModel)
	if err != nil {
		return errors.Wrap(err, "Error while decrypting the model")
	}
	log.Debug("Successfully decrypted model")

	//Decrypt the model inside the TD
	mod := C.CBytes(model)
	C.aimodelbuffer = (*C.char)(mod)
	return nil
}

func UnwrapKey(wrappedKey []byte, pri *rsa.PrivateKey) ([]byte, error) {

	decryptedKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pri, wrappedKey, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error while decrypting the swk")
	}

	return decryptedKey, nil
}

func Decrypt(key, cipherText []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "Error initializing cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating a cipher block")
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("Invalid cipher text")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error decrypting data")
	}

	return plainText, nil
}
