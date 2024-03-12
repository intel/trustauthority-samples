/*
 * Copyright (C) 2024 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/intel/trustauthority-samples/tdxexample/model"
	"github.com/intel/trustauthority-samples/tdxexample/version"
)

type Service interface {
	GetAttestationToken(context.Context) (*GetAttestationTokenResponse, error)
	GetQuote(context.Context, GetQuoteRequest) (*GetQuoteResponse, error)
	GetKey(context.Context, GetKeyRequest) (*GetKeyResponse, error)
	Execute(context.Context, InferRequest) (*InferResponse, error)
	Decrypt(context.Context, GetKeyResponse) (interface{}, error)
	Reset(context.Context) (interface{}, error)
	GetVersion(context.Context) (*version.ServiceVersion, error)
	Provision(context.Context, ProvisionRequest) (interface{}, error)
}

type service struct {
	userData   string
	httpClient *http.Client
	executor   *model.ModelExecutor
}

func NewService(userData string, httpClient *http.Client, executor *model.ModelExecutor) (Service, error) {

	var svc Service
	{
		svc = service{
			userData:   userData,
			httpClient: httpClient,
			executor:   executor,
		}
	}

	svc = LoggingMiddleware()(svc)
	return svc, nil
}

type HandledError struct {
	Code    int
	Message string
}

func (e HandledError) StatusCode() int {
	return e.Code
}

func (e HandledError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
