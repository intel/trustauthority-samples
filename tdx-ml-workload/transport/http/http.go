/*
 * Copyright (C) 2024 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package http

import (
	"context"
	"encoding/json"
	"net/http"

	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/intel/trustauthority-samples/tdxexample/service"
	log "github.com/sirupsen/logrus"
)

const (
	HTTPHeaderKeyContentType       = "Content-Type"
	HTTPHeaderValueApplicationJson = "application/json"
	HTTPHeaderTypeXApiKey          = "x-api-key"
	HTTPHeaderKeyAccept            = "Accept"
	HTTPHeaderKeyAttestationType   = "Attestation-Type"
)

func NewHTTPHandler(svc service.Service) (http.Handler, error) {
	r := mux.NewRouter()
	r.SkipClean(true)

	options := []httpTransport.ServerOption{
		httpTransport.ServerErrorEncoder(errorEncoder),
	}

	{
		prefix := r.PathPrefix("/taa/v1")
		sr := prefix.Subrouter()

		myHandlers := []func(service.Service, *mux.Router, []httpTransport.ServerOption) error{
			setGetVersionHandler,
			setQuoteHandler,
			setModelHandler,
			setKeyHandler,
			setAttestationTokenHandler,
			setProvisionHandler,
		}

		for _, handler := range myHandlers {
			if err := handler(svc, sr, options); err != nil {
				return nil, err
			}
		}
	}

	h := handlers.RecoveryHandler(
		handlers.RecoveryLogger(log.StandardLogger()),
		handlers.PrintRecoveryStack(true),
	)(
		handlers.CombinedLoggingHandler(
			log.StandardLogger().Writer(),
			r,
		),
	)

	return h, nil
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	if handledError, ok := err.(*service.HandledError); ok {
		w.WriteHeader(handledError.Code)
	} else {
		w.WriteHeader(errToCode(err))
	}
	if err := json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()}); err != nil {
		log.WithError(err).Error("Failed to encode error")
	}
}

func errToCode(err error) int {
	switch err {
	case ErrInvalidRequest, ErrJsonDecodeFailed, ErrEmptyRequestBody, ErrTooManyQueryParams, ErrInvalidQueryParam, ErrInvalidFilterCriteria, ErrBase64DecodeFailed:
		return http.StatusBadRequest
	case ErrInvalidContentTypeHeader, ErrInvalidAcceptHeader:
		return http.StatusUnsupportedMediaType
	}
	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}
