/*
 * Copyright (c) 2024 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package version

import (
	"fmt"
)

var Version = ""
var GitHash = ""
var BuildDate = ""

type ServiceVersion struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	GitHash   string `json:"gitHash"`
	BuildDate string `json:"buildDate"`
}

var serviceVersion = ServiceVersion{
	Name:      "TrustAuthority Demo App",
	Version:   Version,
	GitHash:   GitHash,
	BuildDate: BuildDate,
}

func GetVersion() *ServiceVersion {
	return &serviceVersion
}

func (ver *ServiceVersion) String() string {
	return fmt.Sprintf("%s %s-%s [%s]", ver.Name, ver.Version, ver.GitHash, ver.BuildDate)
}
