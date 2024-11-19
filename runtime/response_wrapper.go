package runtime

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type ResultCodeWrapper interface {
	WriteResponse(ctx context.Context, rw http.ResponseWriter, code ResponseCode, msg string, data interface{})
	cannotImplementOutside()
}

type codeInBodyWrapper struct {
	serverErrorCode    ResponseCode
	serverErrorDesc    string
	handlerFacilitator HandlerFacilitator
}

func NewCodeInBodyWrapper(
	hf HandlerFacilitator, serverErrorCode ResponseCode, serverErrorDesc string,
) ResultCodeWrapper {
	return &codeInBodyWrapper{
		serverErrorCode:    serverErrorCode,
		serverErrorDesc:    serverErrorDesc,
		handlerFacilitator: hf,
	}
}

type response struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
}

func (cib *codeInBodyWrapper) WriteResponse(ctx context.Context, rw http.ResponseWriter, code ResponseCode, msg string, data interface{}) {
	respData, marshalError := cib.handlerFacilitator.EncodeJSON(&response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
	if marshalError != nil {
		cib.handlerFacilitator.ErrorHook(ctx, APIError{
			Code:        cib.serverErrorCode,
			Message:     cib.serverErrorDesc,
			SourceError: marshalError,
		})
		respData = []byte(fmt.Sprintf(`{"code":%d,"message":"%s"}`,
			cib.serverErrorCode, cib.serverErrorDesc))
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := rw.Write(respData)
	if err != nil {
		cib.handlerFacilitator.ErrorHook(ctx, APIError{
			Code:        cib.serverErrorCode,
			Message:     cib.serverErrorDesc,
			SourceError: err,
		})
		return
	}
	if marshalError == nil {
		cib.handlerFacilitator.ResultHook(ctx, data)
	}
}

func (cib *codeInBodyWrapper) cannotImplementOutside() {}

type codeInHeaderWrapper struct {
	serverErrorCode ResponseCode
	serverErrorDesc string
	hf              HandlerFacilitator
}

func NewCodeInHeaderWrapper(
	hf HandlerFacilitator, serverErrorCode ResponseCode, serverErrorDesc string,
) ResultCodeWrapper {
	return &codeInHeaderWrapper{
		serverErrorCode: serverErrorCode,
		serverErrorDesc: serverErrorDesc,
		hf:              hf,
	}
}

func (cih *codeInHeaderWrapper) WriteResponse(ctx context.Context, rw http.ResponseWriter, code ResponseCode, msg string, data interface{}) {
	var (
		respData     []byte
		marshalError error
	)
	if data != nil {
		respData, marshalError = cih.hf.EncodeJSON(data)
	}
	if marshalError != nil {
		cih.hf.ErrorHook(ctx, APIError{
			Code:        cih.serverErrorCode,
			Message:     cih.serverErrorDesc,
			SourceError: marshalError,
		})
		rw.Header().Set("X-Hapi-Code", strconv.FormatInt(int64(cih.serverErrorCode), 10))
		rw.Header().Set("X-Hapi-Message", cih.serverErrorDesc)
	} else {
		rw.Header().Set("X-Hapi-Code", strconv.FormatInt(int64(code), 10))
		rw.Header().Set("X-Hapi-Message", msg)
		if respData != nil {
			rw.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, err := rw.Write(respData)
			if err != nil {
				cih.hf.ErrorHook(ctx, APIError{
					Code:        cih.serverErrorCode,
					Message:     cih.serverErrorDesc,
					SourceError: err,
				})
				return
			}
		}
	}
	cih.hf.ResultHook(ctx, data)
}

func (cih *codeInHeaderWrapper) cannotImplementOutside() {}
