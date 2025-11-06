// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"fmt"
	"net/http"
)

// ErrorCode represents the service error code string returned by VikingDB.
type ErrorCode string

// Predefined SDK error codes.
const (
	// HTTP layer errors.
	ErrCodeHTTPRequestFailed ErrorCode = "HTTPRequestFailed"

	// Generic errors.
	ErrCodeUnknown              ErrorCode = "Unknown"
	ErrCodeInvalidParameter     ErrorCode = "InvalidParameter"
	ErrCodeServiceUnavailable   ErrorCode = "ServiceUnavailable"
	ErrCodeTimeout              ErrorCode = "Timeout"
	ErrCodeRequestLimitExceeded ErrorCode = "RequestLimitExceeded"
	ErrCodeUnauthorized         ErrorCode = "Unauthorized"
	ErrCodeForbidden            ErrorCode = "Forbidden"
	ErrCodeNotFound             ErrorCode = "NotFound"

	// Collection related errors.
	ErrCodeCollectionNotExists     ErrorCode = "CollectionNotExists"
	ErrCodeCollectionAlreadyExists ErrorCode = "CollectionAlreadyExists"
	ErrCodeCollectionCreateFailed  ErrorCode = "CollectionCreateFailed"
	ErrCodeCollectionUpdateFailed  ErrorCode = "CollectionUpdateFailed"
	ErrCodeCollectionDeleteFailed  ErrorCode = "CollectionDeleteFailed"

	// Data related errors.
	ErrCodeDataInsertFailed ErrorCode = "DataInsertFailed"
	ErrCodeDataUpdateFailed ErrorCode = "DataUpdateFailed"
	ErrCodeDataDeleteFailed ErrorCode = "DataDeleteFailed"
	ErrCodeDataNotFound     ErrorCode = "DataNotFound"

	// Search related errors.
	ErrCodeSearchFailed   ErrorCode = "SearchFailed"
	ErrCodeIndexNotExists ErrorCode = "IndexNotExists"

	// Embedding related errors.
	ErrCodeEmbeddingFailed ErrorCode = "EmbeddingFailed"
	ErrCodeModelNotFound   ErrorCode = "ModelNotFound"
)

// Error wraps a VikingDB failure with HTTP and internal metadata.
type Error struct {
	// Code is the VikingDB error code string.
	Code ErrorCode `json:"code"`

	// Message describes the failure.
	Message string `json:"message"`

	// StatusCode is the HTTP status returned by the service.
	StatusCode int `json:"status_code,omitempty"`

	// RequestID echoes the server-side request identifier.
	RequestID string `json:"request_id,omitempty"`

	// Err contains the underlying error when available.
	Err error `json:"-"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("vikingdb error: code=%s, message=%s, status_code=%d, err=%v, request_id=%s", e.Code, e.Message, e.StatusCode, e.Err, e.RequestID)
	}
	return fmt.Sprintf("vikingdb error: code=%s, message=%s, status_code=%d, err=%v", e.Code, e.Message, e.StatusCode, e.Err)
}

// Unwrap returns the wrapped error for errors.Is compatibility.
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError constructs an Error with the supplied code and message.
func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewErrorWithStatusCode builds an Error with an explicit HTTP status code.
func NewErrorWithStatusCode(code ErrorCode, message string, statusCode int) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewErrorWithRequestID builds an Error carrying a request identifier.
func NewErrorWithRequestID(code ErrorCode, message string, requestID string, statusCode int) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		RequestID:  requestID,
	}
}

// NewErrorWithCause builds an Error that wraps a lower-level failure.
func NewErrorWithCause(code ErrorCode, message string, cause error, statusCode int) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        cause,
	}
}

// IsRetryableError reports whether the error should be retried.
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	sdkErr, ok := err.(*Error)
	if !ok {
		return false
	}

	switch sdkErr.StatusCode {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}

	switch sdkErr.Code {
	case ErrCodeServiceUnavailable, ErrCodeTimeout, ErrCodeRequestLimitExceeded:
		return true
	}

	return false
}

// NewInvalidParameterError returns a BadRequest error.
func NewInvalidParameterError(message string) *Error {
	return NewErrorWithStatusCode(ErrCodeInvalidParameter, message, http.StatusBadRequest)
}

// NewUnauthorizedError returns an Unauthorized error.
func NewUnauthorizedError(message string) *Error {
	return NewErrorWithStatusCode(ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

// NewForbiddenError returns a Forbidden error.
func NewForbiddenError(message string) *Error {
	return NewErrorWithStatusCode(ErrCodeForbidden, message, http.StatusForbidden)
}

// NewNotFoundError returns a NotFound error.
func NewNotFoundError(message string) *Error {
	return NewErrorWithStatusCode(ErrCodeNotFound, message, http.StatusNotFound)
}

// NewServiceUnavailableError returns a ServiceUnavailable error.
func NewServiceUnavailableError(message string) *Error {
	return NewErrorWithStatusCode(ErrCodeServiceUnavailable, message, http.StatusServiceUnavailable)
}

// NewTimeoutError returns a GatewayTimeout error.
func NewTimeoutError(message string) *Error {
	return NewErrorWithStatusCode(ErrCodeTimeout, message, http.StatusGatewayTimeout)
}

// NewRequestLimitExceededError returns a TooManyRequests error.
func NewRequestLimitExceededError(message string) *Error {
	return NewErrorWithStatusCode(ErrCodeRequestLimitExceeded, message, http.StatusTooManyRequests)
}
