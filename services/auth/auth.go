package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/storage"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var _ Authenticator = (*Auth)(nil)

const (
	salt             = "nJkksjjdxszx120_dssd!xc"
	contextKeyUserID = "user-id"
)

//  Auth implements Authenticator interface methods for user authorisation.
type Auth struct {
	tokenAuth *jwtauth.JWTAuth
	storage   storage.Storage
}

//  NewAuth init new Auth.
func NewAuth(s storage.Storage) (*Auth, error) {
	return &Auth{
		tokenAuth: jwtauth.New("HS256", []byte("secret"), nil),
		storage:   s,
	}, nil
}

//  CreateUser creates new user.
//  If user exist, returns ErrorUserAlreadyExist.
func (a *Auth) CreateUser(ctx context.Context, login string, password string) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), 10)
	l := len(hash)
	fmt.Printf("length:%v", l)
	if err != nil {
		return model.User{}, fmt.Errorf("ошибка добавления пользователя: %w", err)
	}

	user := model.User{
		Login:        login,
		PasswordHash: string(hash),
	}

	user, err = a.storage.User().Create(ctx, user)
	if err != nil {
		if errors.Is(err, model.ErrorConflictSaveUser) {
			return model.User{}, ErrorUserAlreadyExist
		}

		return model.User{}, fmt.Errorf("ошибка добавления пользователя: %w", err)
	}

	return user, err
}

//  Authenticate checks user login, password and return.
//  If user not founded, or wrong password, returns ErrorWrongAuthData.
func (a *Auth) Authenticate(ctx context.Context, login string, password string) (model.User, error) {
	user, err := a.storage.User().GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrorItemNotFound) {
			return model.User{}, ErrorWrongAuthData
		}

		return model.User{}, fmt.Errorf("ошибка авторизации пользователя: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password+salt)); err != nil {
		return model.User{}, ErrorWrongAuthData
	}

	return user, nil
}

//  EncodeTokenUserID encodes token with user_id claim.
func (a *Auth) EncodeTokenUserID(userID uuid.UUID) (string, error) {
	_, tokenString, err := a.tokenAuth.Encode(map[string]interface{}{"user_id": userID.String()})
	if err != nil {
		return "", fmt.Errorf("ошибка генерации токена для пользователя: %w", err)
	}

	return tokenString, nil
}

//  TokenAuth returns JWTAuth pointer
func (a *Auth) TokenAuth() *jwtauth.JWTAuth {
	return a.tokenAuth
}
