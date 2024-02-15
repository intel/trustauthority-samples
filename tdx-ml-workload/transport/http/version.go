/*
 * Copyright (C) 2022 Intel Corporation
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

func setGetVersionHandler(svc service.Service, router *mux.Router, options []httpTransport.ServerOption) error {
	getVersionHandler := httpTransport.NewServer(
		makeGetVersionHTTPEndpoint(svc),
		httpTransport.NopRequestDecoder,
		httpTransport.EncodeJSONResponse,
		options...,
	)

	router.Handle("/version", getVersionHandler).Methods(http.MethodGet)
	return nil
}

func makeGetVersionHTTPEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		return svc.GetVersion(ctx)
	}
}
