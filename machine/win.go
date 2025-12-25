//go:build windows

// Copyright (C) 2025 Murilo Gomes Julio
// SPDX-License-Identifier: MIT

// Site: https://www.mugomes.com.br

package machine

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

type winProvider struct{}

func init() {
	Current = winProvider{}
}

func (winProvider) GetMachineID() string {
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Cryptography`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return "unknown"
	}
	defer k.Close()

	id, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(id)
}
