/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package service

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type InferRequest struct {
	Pregnancies   float32 `json:"pregnancies,string"`
	BloodGlucose  float32 `json:"blood-glucose,string"`
	BloodPressure float32 `json:"blood-pressure,string"`
	SkinThickness float32 `json:"skin-thickness,string"`
	Insulin       float32 `json:"insulin,string"`
	BMI           float32 `json:"bmi,string"`
	Age           float32 `json:"age,string"`
	DBF           float32 `json:"dbf,string"`
}

type InferResponse struct {
	HighRisk int `json:"high-risk"`
}

func (t *InferResponse) Headers() http.Header {
	return corsHeaders
}

type ModelResponse struct {
	statusCode int
}

func (t *ModelResponse) Headers() http.Header {
	return corsHeaders
}

func (d *ModelResponse) StatusCode() int {
	return d.statusCode
}

func (mw loggingMiddleware) Decrypt(ctx context.Context, req GetKeyResponse) (interface{}, error) {
	var err error
	defer func(begin time.Time) {
		log.Tracef("Decrypt took %s since %s", time.Since(begin), begin)
		if err != nil {
			log.WithError(err)
		}
	}(time.Now())
	resp, err := mw.next.Decrypt(ctx, req)
	return resp, err
}

func (svc service) Decrypt(_ context.Context, req GetKeyResponse) (interface{}, error) {

	err := svc.executor.DecryptModel(req.WrappedSwk, req.WrappedKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not decrypt model")
	}
	return &ModelResponse{http.StatusNoContent}, nil
}

func (mw loggingMiddleware) Execute(ctx context.Context, req InferRequest) (*InferResponse, error) {
	var err error
	defer func(begin time.Time) {
		log.Tracef("Execute took %s since %s", time.Since(begin), begin)
		if err != nil {
			log.WithError(err)
		}
	}(time.Now())
	resp, err := mw.next.Execute(ctx, req)
	return resp, err
}

func (svc service) Execute(_ context.Context, req InferRequest) (*InferResponse, error) {

	res, err := svc.executor.ExecuteModel(req.Pregnancies, req.BloodGlucose,
		req.BloodPressure, req.SkinThickness,
		req.Insulin, req.BMI, req.DBF, req.Age)
	if err != nil {
		return nil, errors.Wrap(err, "could not execute model")
	}

	resp := &InferResponse{
		HighRisk: res,
	}
	return resp, nil
}

func (mw loggingMiddleware) Reset(ctx context.Context) (interface{}, error) {
	var err error
	defer func(begin time.Time) {
		log.Tracef("Reset took %s since %s", time.Since(begin), begin)
		if err != nil {
			log.WithError(err)
		}
	}(time.Now())
	resp, err := mw.next.Reset(ctx)
	return resp, err
}

func (svc service) Reset(_ context.Context) (interface{}, error) {

	err := svc.executor.ResetModel()
	if err != nil {
		return nil, errors.Wrap(err, "could not reset model")
	}
	return &ModelResponse{http.StatusNoContent}, nil
}
