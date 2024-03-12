/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package service

import (
	"context"
	"encoding/base64"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type GetQuoteRequest struct {
	Nonce []byte
}

type GetQuoteResponse struct {
	Quote    []byte
	UserData []byte
}

func (mw loggingMiddleware) GetQuote(ctx context.Context, req GetQuoteRequest) (*GetQuoteResponse, error) {
	var err error
	defer func(begin time.Time) {
		log.Tracef("GetQuote took %s since %s", time.Since(begin), begin)
		if err != nil {
			log.WithError(err)
		}
	}(time.Now())
	resp, err := mw.next.GetQuote(ctx, req)
	return resp, err
}

func (svc service) GetQuote(_ context.Context, req GetQuoteRequest) (*GetQuoteResponse, error) {

	quote, userData, err := collectEvidence(base64.StdEncoding.EncodeToString(req.Nonce), svc.userData)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch quote")
	}

	resp := &GetQuoteResponse{
		Quote:    quote,
		UserData: userData,
	}
	return resp, nil
}

func collectEvidence(nonce, userData string) ([]byte, []byte, error) {

	out, err := exec.Command(CLI, "quote", "--nonce", nonce, "--user-data", userData).Output()
	if err != nil {
		return nil, nil, err
	}
	runtimeData, err := base64.StdEncoding.DecodeString(userData)
	if err != nil {
		return nil, nil, err
	}

	return out[:], runtimeData, nil
}
