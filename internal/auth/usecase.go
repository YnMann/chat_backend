package auth

import (
	"context"

	"github.com/YnMann/chat_backend/internal/models"
)

const CtxUserKey = "auth"

type UseCase interface {
	SignUp(ctx context.Context, data *models.User) error
	SignIn(ctx context.Context, username, password string) (string, error)
	ParseToken(ctx context.Context, accessToken string) (*models.User, error)
}
