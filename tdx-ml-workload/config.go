/*
 * Copyright (c) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package main

import (
	"encoding/base64"
	"net/url"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	envServicePort              = "SERVICE_PORT"
	envEnableLogCaller          = "LOG_CALLER"
	envLogLevel                 = "LOG_LEVEL"
	envSkipTlsVerification      = "SKIP_TLS_VERIFICATION"
	envHttpClientTimeoutSeconds = "HTTP_CLIENT_TIMEOUT_IN_SECONDS"

	envITAUrl    = "ITA_API_URL"
	envITAApiKey = "ITA_API_KEY"

	defaultPort        = "12780"
	defaultLogLevel    = "info"
	defaultHttpTimeout = "10"
)

type Configuration struct {
	Port                int
	LogCaller           bool
	LogLevel            log.Level
	SkipTLSVerification bool
	HTTPClientTimeout   int

	ITAURL    string
	ITAAPIKey string
}

func configure() (*Configuration, error) {
	log.SetFormatter(&log.JSONFormatter{})
	c, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	log.SetReportCaller(c.LogCaller)
	log.SetLevel(c.LogLevel)
	return c, nil
}

func NewConfigFromEnv() (*Configuration, error) {

	// Setup defaults by associating the structure's field names to the defaults.
	viper.SetDefault("Port", defaultPort)
	viper.SetDefault("LogCaller", "false")
	viper.SetDefault("SkipTlsVerification", "true")
	viper.SetDefault("HTTPClientTimeout", defaultHttpTimeout)

	// map structure field names to env var names (log level is handled manually below)
	envBinding := map[string]string{
		"Port":                envServicePort,
		"LogCaller":           envEnableLogCaller,
		"SkipTLSVerification": envSkipTlsVerification,
		"HTTPClientTimeout":   envHttpClientTimeoutSeconds,
		"ITAUrl":              envITAUrl,
		"ITAApiKey":           envITAApiKey,
	}

	for fieldName, envVar := range envBinding {
		err := viper.BindEnv(fieldName, envVar)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to bind env var %s to field %s", envVar, fieldName)
		}
	}

	var conf Configuration
	err := viper.Unmarshal(&conf)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse config from environment variables")
	}

	envLevel := os.Getenv(envLogLevel)
	if envLevel == "" {
		conf.LogLevel = log.InfoLevel
	} else if logLevel, err := log.ParseLevel(envLevel); err != nil {
		log.Warnf("Failed to parse log level %q, defaulting to 'info' level", envLevel)
		conf.LogLevel = log.InfoLevel
	} else {
		conf.LogLevel = logLevel
	}

	log.WithFields(log.Fields{
		"Port":                conf.Port,
		"LogLevel":            conf.LogLevel,
		"LogCaller":           conf.LogCaller,
		"SkipTLSVerification": conf.SkipTLSVerification,
		"HTTPClientTimeout":   conf.HTTPClientTimeout,
		"ITAUrl":              conf.ITAURL,
	}).Info("Parse configs from environment")

	return &conf, nil
}

func (conf *Configuration) Validate() error {

	if conf.Port < 1024 || conf.Port > 65535 {
		return errors.New("Configured port is not valid")
	}

	if conf.ITAURL == "" || conf.ITAAPIKey == "" {
		return errors.New("Either ITA API URL or APIKey is missing")
	}

	_, err := url.Parse(conf.ITAURL)
	if err != nil {
		return errors.Wrap(err, "ITA API URL is not a valid url")
	}

	_, err = base64.StdEncoding.DecodeString(conf.ITAAPIKey)
	if err != nil {
		return errors.Wrap(err, "ITA ApiKey is not a valid base64 string")
	}

	return nil
}