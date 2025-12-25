//go:build darwin

// Copyright (C) 2025 Murilo Gomes Julio
// SPDX-License-Identifier: MIT

// Site: https://www.mugomes.com.br

package machine

import (
	"os/exec"
	"regexp"
)

type macProvider struct{}

func init() {
	Current = macProvider{}
}

func (macProvider) GetMachineID() string {
	out, err := exec.Command(
		"ioreg", "-rd1", "-c", "IOPlatformExpertDevice",
	).Output()
	if err != nil {
		return "unknown"
	}

	re := regexp.MustCompile(`"IOPlatformUUID" = "([^"]+)"`)
	m := re.FindSubmatch(out)
	if len(m) < 2 {
		return "unknown"
	}

	return string(m[1])
}
