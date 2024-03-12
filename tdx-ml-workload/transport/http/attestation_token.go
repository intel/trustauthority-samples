/*
 * Copyright (C) 2024 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package http

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/intel/trustauthority-samples/tdxexample/service"
)

func setAttestationTokenHandler(svc service.Service, router *mux.Router, options []httpTransport.ServerOption) error {

	getAttestationTokenHandler := httpTransport.NewServer(
		makeGetAttestationTokenHTTPEndpoint(svc),
		httpTransport.NopRequestDecoder,
		httpTransport.EncodeJSONResponse,
		options...,
	)

	router.Handle("/token", getAttestationTokenHandler).Methods(http.MethodGet)
	router.Handle("/token", optionsHandler).Methods(http.MethodOptions)

	return nil
}

func makeGetAttestationTokenHTTPEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.GetAttestationToken(ctx)
	}
}
