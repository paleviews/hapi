package runtime

import (
	"context"
	"net/http"
)

type HandlerFacilitator interface {
	VerifyAuthToken(context.Context, string) (context.Context, error)
	EncodeJSON(interface{}) ([]byte, error)
	DecodeJSON([]byte, interface{}) error
	GetPathParam(_ *http.Request, key string) string
	ResultHook(_ context.Context, data interface{})
	ErrorHook(context.Context, APIError)
}
