/*
 * Copyright (C) 2024 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package service

import (
	"bytes"
	"context"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	CLI = "trustauthority-cli"
)

type GetAttestationTokenResponse struct {
	AttestationToken string `json:"attestation_token"`
}

func (t *GetAttestationTokenResponse) Headers() http.Header {
	return corsHeaders
}

func (mw loggingMiddleware) GetAttestationToken(ctx context.Context) (*GetAttestationTokenResponse, error) {
	var err error
	defer func(begin time.Time) {
		log.Tracef("GetAttestationToken took %s since %s", time.Since(begin), begin)
		if err != nil {
			log.WithError(err)
		}
	}(time.Now())
	resp, err := mw.next.GetAttestationToken(ctx)
	return resp, err
}

func (svc service) GetAttestationToken(_ context.Context) (*GetAttestationTokenResponse, error) {

	var policyIds string
	cmd := exec.Command(CLI, "token", "--config", "config.json", "--user-data", svc.userData, "--policy-ids", policyIds, "--no-eventlog")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, errors.Wrapf(err, "could not fetch token: %v", string(stderr.Bytes()))
	}

	resp := &GetAttestationTokenResponse{
		AttestationToken: strings.TrimSpace(string(stdout.Bytes())),
	}
	return resp, nil
}
