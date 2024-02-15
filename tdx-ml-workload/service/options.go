/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package service

import (
	"net/http"
)

var corsHeaders = map[string][]string{
	"Access-Control-Allow-Origin":  {"*"},
	"Access-Control-Allow-Methods": {"POST, GET, OPTIONS, PUT, DELETE"},
	"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Allow-Origin"},
}

type OptionsResponse struct {
}

func (t *OptionsResponse) Headers() http.Header {
	return corsHeaders
}
