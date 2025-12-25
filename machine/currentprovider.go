// Copyright (C) 2025 Murilo Gomes Julio
// SPDX-License-Identifier: MIT

// Site: https://www.mugomes.com.br

package machine

type Provider interface {
	GetMachineID() string
}

var Current Provider