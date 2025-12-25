//go:build linux

// Copyright (C) 2025 Murilo Gomes Julio
// SPDX-License-Identifier: MIT

// Site: https://www.mugomes.com.br

package machine

import (
	"os"
	"strings"
)

type linuxProvider struct{}

func init() {
	Current = linuxProvider{}
}

func (linuxProvider) GetMachineID() string {
	b, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(b))
}
