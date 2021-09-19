package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/paleviews/hapi/example/todo/auth"
	"github.com/paleviews/hapi/runtime"
)

type Facilitator struct {
	mux    *chi.Mux
	logger *zap.Logger
}

func NewFacilitator() *Facilitator {
	logger, _ := zap.NewDevelopment()
	return &Facilitator{
		mux:    chi.NewMux(),
		logger: logger,
	}
}

func (f *Facilitator) VerifyAuthToken(ctx context.Context, token string) (context.Context, error) {
	if token == "super_secret" {
		return auth.NewCtxWithUserID(ctx, 2022), nil
	}
	return nil, errors.New("invalid token")
}

func (f *Facilitator) EncodeJSON(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (f *Facilitator) DecodeJSON(b []byte, dst interface{}) error {
	return json.Unmarshal(b, dst)
}

func (f *Facilitator) GetPathParam(req *http.Request, key string) string {
	return chi.URLParam(req, key)
}

func (f *Facilitator) ResultHook(context.Context, interface{}) {}

func (f *Facilitator) ErrorHook(_ context.Context, err runtime.APIError) {
	f.logger.Error("error hook", zap.Any("err", &err))
}

func (f *Facilitator) Mux() *chi.Mux {
	return f.mux
}
