/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/intel/trustauthority-samples/tdxexample/service"
	log "github.com/sirupsen/logrus"
)

func setModelHandler(svc service.Service, router *mux.Router, options []httpTransport.ServerOption) error {

	decryptHandler := httpTransport.NewServer(
		makeDecryptHTTPEndpoint(svc),
		decodeDecryptHTTPRequest,
		httpTransport.EncodeJSONResponse,
		options...,
	)

	router.Handle("/decrypt", decryptHandler).Methods(http.MethodPost)
	router.Handle("/decrypt", optionsHandler).Methods(http.MethodOptions)

	executeHandler := httpTransport.NewServer(
		makeExecuteHTTPEndpoint(svc),
		decodeExecuteHTTPRequest,
		httpTransport.EncodeJSONResponse,
		options...,
	)

	router.Handle("/execute", executeHandler).Methods(http.MethodPost)
	router.Handle("/execute", optionsHandler).Methods(http.MethodOptions)

	resetHandler := httpTransport.NewServer(
		makeResetHTTPEndpoint(svc),
		httpTransport.NopRequestDecoder,
		httpTransport.EncodeJSONResponse,
		options...,
	)

	router.Handle("/reset", resetHandler).Methods(http.MethodPost)
	router.Handle("/reset", optionsHandler).Methods(http.MethodOptions)

	return nil
}

func makeDecryptHTTPEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.GetKeyResponse)
		return svc.Decrypt(ctx, req)
	}
}

func makeExecuteHTTPEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.InferRequest)
		return svc.Execute(ctx, req)
	}
}

func makeResetHTTPEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.Reset(ctx)
	}
}

func decodeDecryptHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {

	if r.Header.Get(HTTPHeaderKeyContentType) != HTTPHeaderValueApplicationJson {
		log.Error(ErrInvalidContentTypeHeader.Error())
		return nil, ErrInvalidContentTypeHeader
	}

	if r.ContentLength == 0 {
		log.Error(ErrEmptyRequestBody.Error())
		return nil, ErrEmptyRequestBody
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req service.GetKeyResponse
	err := dec.Decode(&req)
	if err != nil {
		log.WithError(err).Error(ErrJsonDecodeFailed.Error())
		return nil, ErrJsonDecodeFailed
	}

	return req, nil
}

func decodeExecuteHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {

	if r.Header.Get(HTTPHeaderKeyContentType) != HTTPHeaderValueApplicationJson {
		log.Error(ErrInvalidContentTypeHeader.Error())
		return nil, ErrInvalidContentTypeHeader
	}

	if r.ContentLength == 0 {
		log.Error(ErrEmptyRequestBody.Error())
		return nil, ErrEmptyRequestBody
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req service.InferRequest
	err := dec.Decode(&req)
	if err != nil {
		log.WithError(err).Error(ErrJsonDecodeFailed.Error())
		return nil, ErrJsonDecodeFailed
	}

	return req, nil
}
