package auth

import (
	"context"
	"github.com/atrush/diploma.git/model"
	"github.com/google/uuid"
)

//  Authenticator is the interface that wraps methods user identification, authentication, authorisation.
type Authenticator interface {
	//TokenAuth() *jwtauth.JWTAuth
	CreateUser(ctx context.Context, login string, password string) (model.User, error)
	Authenticate(ctx context.Context, login string, password string) (model.User, error)
	EncodeTokenUserID(userID uuid.UUID) (string, error)
}
