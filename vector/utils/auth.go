// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"net/http"

	"github.com/volcengine/volc-sdk-golang/base"
)

const (
	defaultRegion  = "cn-beijing"
	defaultService = "vikingdb"
)

// SignRequest keeps backward compatibility with the previous signature and signs using the default region.
func SignRequest(req *http.Request, ak, sk string) *http.Request {
	return SignRequestWithRegion(req, ak, sk, defaultRegion)
}

// SignRequestWithRegion signs the HTTP request with the provided credentials and region.
func SignRequestWithRegion(req *http.Request, ak, sk, region string) *http.Request {
	credential := base.Credentials{
		AccessKeyID:     ak,
		SecretAccessKey: sk,
		Service:         defaultService,
		Region:          region,
	}
	return credential.Sign(req)
}
