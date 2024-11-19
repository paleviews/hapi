// Code generated by proto-gen-hapi. DO NOT EDIT.
// version:
// - protoc-gen-hapi v0.1.0
// - protoc          v5.28.3
// source: codes/response_code.proto

package codes

import (
	runtime "github.com/paleviews/hapi/runtime"
)

type ResponseCode = runtime.ResponseCode

const (
	// just ok
	ResponseCode_RESPONSE_CODE_OK              ResponseCode = 0
	ResponseCode_RESPONSE_CODE_INVALID_INPUT   ResponseCode = 1
	ResponseCode_RESPONSE_CODE_UNAUTHENTICATED ResponseCode = 2
	ResponseCode_RESPONSE_CODE_NOT_FOUND       ResponseCode = 3
	ResponseCode_RESPONSE_CODE_SERVER_ERROR    ResponseCode = 99
	ResponseCode_RESPONSE_CODE_UNIMPLEMENTED   ResponseCode = 1001
)

func GetDescByResponseCode(code ResponseCode) string {
	switch code {
	case ResponseCode_RESPONSE_CODE_OK:
		return "ok"
	case ResponseCode_RESPONSE_CODE_INVALID_INPUT:
		return "invalid_input"
	case ResponseCode_RESPONSE_CODE_UNAUTHENTICATED:
		return "unauthenticated"
	case ResponseCode_RESPONSE_CODE_NOT_FOUND:
		return "not_found"
	case ResponseCode_RESPONSE_CODE_SERVER_ERROR:
		return "server_error"
	case ResponseCode_RESPONSE_CODE_UNIMPLEMENTED:
		return "unimplemented"
	default:
		return ""
	}
}

func APIErrorFromResponseCode(code ResponseCode, src error) runtime.APIError {
	return runtime.NewAPIError(code, GetDescByResponseCode(code), src)
}
