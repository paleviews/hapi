package runtime

import (
	"errors"
	"net/http"
)

func GetBearerTokenInHeader(req *http.Request) (string, error) {
	s := req.Header.Get("Authorization")
	const (
		prefix    = "Bearer "
		prefixLen = len(prefix)
	)
	if len(s) < prefixLen {
		return "", errors.New("empty authorization in header")
	}
	if s[:prefixLen] != prefix {
		return "", errors.New("invalid bearer authorization in header")
	}
	return s[prefixLen:], nil
}

func AuthMiddleware(
	hf HandlerFacilitator, tokenGetter func(*http.Request) (string, error),
	unauthenticatedCode ResponseCode, unauthenticatedDesc string, responseWrapper ResultCodeWrapper,
) func(func(rw http.ResponseWriter, req *http.Request)) func(rw http.ResponseWriter, req *http.Request) {

	return func(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
		return func(rw http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			token, err := tokenGetter(req)
			if err != nil {
				hf.ErrorHook(ctx, NewAPIError(unauthenticatedCode, unauthenticatedDesc, err))
				responseWrapper.WriteResponse(ctx, rw, unauthenticatedCode, unauthenticatedDesc, nil)
				return
			}
			newCtx, err := hf.VerifyAuthToken(ctx, token)
			if err != nil {
				hf.ErrorHook(ctx, NewAPIError(unauthenticatedCode, unauthenticatedDesc, err))
				responseWrapper.WriteResponse(ctx, rw, unauthenticatedCode, unauthenticatedDesc, nil)
				return
			}
			next(rw, req.WithContext(newCtx))
		}
	}
}
