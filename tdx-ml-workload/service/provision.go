/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package service

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ProvisionRequest struct {
	ApiKey string `json:"api_key"`
}

func (mw loggingMiddleware) Provision(ctx context.Context, req ProvisionRequest) (interface{}, error) {
	var err error
	defer func(begin time.Time) {
		log.Tracef("Provision took %s since %s", time.Since(begin), begin)
		if err != nil {
			log.WithError(err)
		}
	}(time.Now())
	resp, err := mw.next.Provision(ctx, req)
	return resp, err
}

func (svc service) Provision(_ context.Context, req ProvisionRequest) (interface{}, error) {
	_, err := base64.URLEncoding.DecodeString(req.ApiKey)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Api key, must be base64 string")
	}

	return &ModelResponse{http.StatusNoContent}, nil
}
