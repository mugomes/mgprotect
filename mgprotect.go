// Copyright (C) 2025 Murilo Gomes Julio
// SPDX-License-Identifier: MIT

// Site: https://www.mugomes.com.br

package mgprotect

import (
	"bytes"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/mugomes/mgprotect/machine"
	"os"
	"strings"
)

type licenseFile struct {
	Serial       string
	ProductID    byte
	MajorVersion byte
	MachineHash  []byte
	Signature    []byte
}

type MGPROTECT struct {
	productID    byte
	majorVersion byte
	publicKey    ed25519.PublicKey
	internalKey  []byte
	K            string
	licenseFile  licenseFile
}

func New() *MGPROTECT {
	return &MGPROTECT{}
}

func (mgp *MGPROTECT) SetProductID(value byte) {
	mgp.productID = value
}

func (mgp *MGPROTECT) SetMajorVersion(value byte) {
	mgp.majorVersion = value
}

func (mgp *MGPROTECT) SetPublicKey(value ed25519.PublicKey) {
	mgp.publicKey = value
}

func (mgp *MGPROTECT) SetInternalKey(value []byte) {
	mgp.internalKey = value
}

// Validação do Form
func checksum(b []byte) byte {
	var c byte
	for _, v := range b {
		c ^= v
	}
	return c
}

const (
	ERRO_SERIAL_INVALIDO                = 2
	ERRO_CHECKSUM_INVALIDO              = 3
	ERRO_PRODUTO_INVALIDO               = 4
	ERRO_LICENCA_NAO_VALIDA_PARA_VERSAO = 5
	ERRO_ASSINATURA_INVALIDA            = 6
)

func (mgp *MGPROTECT) Validate(input string) int {
	var sb strings.Builder
	sb.Grow(len(input))
	for _, c := range input {
		if c != '-' && c != ' ' {
			sb.WriteRune(c)
		}
	}

	dec := base64.StdEncoding.WithPadding(base64.NoPadding)
	raw, err := dec.DecodeString(sb.String())
	if err != nil || len(raw) != 4+ed25519.SignatureSize {
		return ERRO_SERIAL_INVALIDO
	}

	payload := raw[:4]
	sig := raw[4:]

	if payload[3] != checksum(payload[:3]) {
		return ERRO_CHECKSUM_INVALIDO
	}

	if payload[0] != mgp.productID {
		return ERRO_PRODUTO_INVALIDO
	}

	if payload[1] != mgp.majorVersion {
		return ERRO_LICENCA_NAO_VALIDA_PARA_VERSAO
	}

	if !ed25519.Verify(mgp.publicKey, payload, sig) {
		return ERRO_ASSINATURA_INVALIDA
	}

	return 001
}

// Validação Boot
func (mgp *MGPROTECT) signLocal(data []byte) []byte {
	h := hmac.New(sha256.New, mgp.internalKey)
	h.Write(data)
	return h.Sum(nil)
}

var mc machine.Provider = machine.Current

func (mgp *MGPROTECT) SetK(value string) {
	mgp.K = value
}

func (mgp *MGPROTECT) SaveLicense(path string, serial string) error {
	mid := mc.GetMachineID()

	h := sha256.Sum256([]byte(mid))
	sSerial, _ := encrypt(serial, mgp.K)
	lf := licenseFile{
		Serial:       sSerial,
		ProductID:    mgp.productID,
		MajorVersion: mgp.majorVersion,
		MachineHash:  h[:],
	}

	raw, _ := json.Marshal(lf)
	lf.Signature = mgp.signLocal(raw)

	final, _ := json.Marshal(lf)
	return os.WriteFile(path, final, 0600)
}

func (mgp *MGPROTECT) LoadAndValidate(path string) bool {
	b, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var lf licenseFile
	if json.Unmarshal(b, &lf) != nil {
		return false
	}

	sig := lf.Signature
	lf.Signature = nil

	raw, _ := json.Marshal(lf)
	if !hmac.Equal(sig, mgp.signLocal(raw)) {
		return false
	}

	sSerial, _ := decrypt(lf.Serial, mgp.K)
	if mgp.Validate(sSerial) != 001 {
		return false
	}

	mid := mc.GetMachineID()
	h := sha256.Sum256([]byte(mid))

	return bytes.Equal(lf.MachineHash, h[:])
}
