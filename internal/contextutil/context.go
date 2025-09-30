package contextutil

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	if userID, ok := ctx.Value(UserIDContextKey).(uuid.UUID); ok {
		return userID, true
	}
	return uuid.Nil, false
}

func SetUserIDInContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}
