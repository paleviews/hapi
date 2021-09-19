package auth

import (
	"context"
)

type ctxKeyUserID struct{}

func NewCtxWithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, ctxKeyUserID{}, userID)
}

func UserIDFromCtx(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(ctxKeyUserID{}).(int64)
	return id, ok
}
