// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
)

// ParseJSONUseNumber decodes input into target while preserving numeric precision via json.Number.
func ParseJSONUseNumber(input []byte, target interface{}) error {
	if target == nil {
		return errors.New("ParseJSONUseNumber: target must not be nil")
	}
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.UseNumber()
	return decoder.Decode(target)
}

// SerializeToJSON marshals the provided value into JSON bytes.
func SerializeToJSON(source interface{}) ([]byte, error) {
	return json.Marshal(source)
}
