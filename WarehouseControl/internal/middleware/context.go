package middleware

import (
	"context"
	"l3/WarehouseControl/internal/models"
)

type contextKey string

const currentUserKey contextKey = "current_user"

func SetCurrentUser(ctx context.Context, user *models.CurrentUser) context.Context {
	return context.WithValue(ctx, currentUserKey, user)
}

func GetCurrentUser(ctx context.Context) (*models.CurrentUser, bool) {
	user, ok := ctx.Value(currentUserKey).(*models.CurrentUser)
	return user, ok
}
