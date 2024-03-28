/*
 * Copyright (c) 2024 Intel Corporation
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
	envSanList                  = "SAN_LIST"
	envSkipTlsVerification      = "SKIP_TLS_VERIFICATION"
	envHttpReadHeaderTimeoutSec = "HTTP_READ_HEADER_TIMEOUT_IN_SECONDS"

	envTrustAuthorityAPIUrl = "TRUSTAUTHORITY_API_URL"
	envTrustAuthorityAPIKey = "TRUSTAUTHORITY_API_KEY"

	defaultSanList     = "127.0.0.1,localhost"
	defaultPort        = "12780"
	defaultLogLevel    = "info"
	defaultHttpTimeout = "10"
)

type Configuration struct {
	Port                int
	SanList             string
	LogCaller           bool
	LogLevel            log.Level
	SkipTLSVerification bool
	HTTPReadHdrTimeout  int

	TrustAuthorityUrl string
	TrustAuthorityKey string
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
	viper.SetDefault("SanList", defaultSanList)
	viper.SetDefault("SkipTlsVerification", "false")
	viper.SetDefault("HTTPReadHdrTimeout", defaultHttpTimeout)

	// map structure field names to env var names (log level is handled manually below)
	envBinding := map[string]string{
		"Port":                envServicePort,
		"SanList":             envSanList,
		"LogCaller":           envEnableLogCaller,
		"SkipTLSVerification": envSkipTlsVerification,
		"HTTPReadHdrTimeout":  envHttpReadHeaderTimeoutSec,
		"TrustAuthorityUrl":   envTrustAuthorityAPIUrl,
		"TrustAuthorityKey":   envTrustAuthorityAPIKey,
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
		"SanList":             conf.SanList,
		"LogLevel":            conf.LogLevel,
		"LogCaller":           conf.LogCaller,
		"SkipTLSVerification": conf.SkipTLSVerification,
		"HTTPReadHdrTimeout":  conf.HTTPReadHdrTimeout,
		"TrustAuthorityUrl":   conf.TrustAuthorityUrl,
	}).Info("Parse configs from environment")

	return &conf, nil
}

func (conf *Configuration) Validate() error {

	if conf.Port < 1024 || conf.Port > 65535 {
		return errors.New("Configured port is not valid")
	}

	if conf.TrustAuthorityUrl == "" || conf.TrustAuthorityKey == "" {
		return errors.New("Either Trust Authority API URL or APIKey is missing")
	}

	_, err := url.Parse(conf.TrustAuthorityUrl)
	if err != nil {
		return errors.Wrap(err, "Trust Authority API URL is not a valid url")
	}

	_, err = base64.StdEncoding.DecodeString(conf.TrustAuthorityKey)
	if err != nil {
		return errors.Wrap(err, "Trust Authority ApiKey is not a valid base64 string")
	}

	return nil
}
