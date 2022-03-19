package api

import (
	"errors"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/services/auth"
	mk "github.com/atrush/diploma.git/services/auth/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"net/http"
	"testing"
)

var (
	mockUser = model.User{
		ID:           uuid.New(),
		Login:        "user_1",
		PasswordHash: "hash",
	}
	reqAuth = "{\"login\": \"iamuser\",\"password\": \"123456\"}"
)

//  TestHandler_Login tests user register handler
func TestHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []TestRoute{
		{
			name:    "return 200 if authenticated",
			method:  http.MethodPost,
			url:     "/api/user/login",
			svcAuth: authOk(ctrl),
			headers: map[string]string{"Content-Type": "application/json"},
			body:    reqAuth,

			expectedHeaders: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer tokentoken",
			},
			expectedCode: 200,
		},
		{
			name:         "return 500 if internal error",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authErrServer(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         reqAuth,
			expectedCode: 500,
		},
		{
			name:         "return 400 if wrong json format, missed quotes",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         "{\"login: \"iamuser\",\"password\": \"123456\"}",
			expectedCode: 400,
		},
		{
			name:         "return 401 if wrong login/password",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authErrWrongPass(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         "{\"login\": \"iamuser\",\"password\": \"123456\"}",
			expectedCode: 401,
		},
		{
			name:         "return 415 if wrong content type",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			body:         reqAuth,
			expectedCode: 415,
		},
		{
			name:         "return 400 if empty login",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         "{\"login: \"\",\"password\": \"123456\"}",
			expectedCode: 400,
		},
		{
			name:         "return 400 if empty password",
			method:       http.MethodPost,
			url:          "/api/user/login",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         "{\"login\": \"iamuser\",\"password\": \"\"}",
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.CheckTest(t)
		})
	}

}

//  TestHandler_Register tests user login handler
func TestHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []TestRoute{
		{
			name:    "return 200 if registered",
			method:  http.MethodPost,
			url:     "/api/user/register",
			svcAuth: registerOk(ctrl),
			headers: map[string]string{"Content-Type": "application/json"},
			body:    reqAuth,

			expectedHeaders: map[string]string{"Content-Type": "application/json", "Authorization": "Bearer tokentoken"},
			expectedCode:    200,
		},
		{
			name:         "return 500 if internal error",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      registerErrServer(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         reqAuth,
			expectedCode: 500,
		},
		{
			name:         "return 409 if user exist",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      registerErrExist(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         reqAuth,
			expectedCode: 409,
		},
		{
			name:         "return 400 if wrong json format, missed quotes",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         "{\"login: \"iamuser\",\"password\": \"123456\"}",
			expectedCode: 400,
		},
		{
			name:         "return 415 if wrong content type",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			body:         reqAuth,
			expectedCode: 415,
		},
		{
			name:         "return 400 if empty login",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         "{\"login: \"\",\"password\": \"123456\"}",
			expectedCode: 400,
		},
		{
			name:         "return 400 if empty password",
			method:       http.MethodPost,
			url:          "/api/user/register",
			svcAuth:      authEmpty(ctrl),
			headers:      map[string]string{"Content-Type": "application/json"},
			body:         "{\"login: \"iamuser\",\"password\": \"\"}",
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.CheckTest(t)
		})
	}

}

/*  Auth mocks  */
func authEmpty(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	return authMock
}

/* Mocks for authenticate handler */
func authOk(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().Authenticate(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockUser, nil)
	authMock.EXPECT().EncodeTokenUserID(gomock.Any()).Return("tokentoken", nil)
	return authMock
}
func authErrWrongPass(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().Authenticate(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.User{}, auth.ErrorWrongAuthData)
	return authMock
}
func authErrServer(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().Authenticate(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.User{}, errors.New("server error"))
	return authMock
}

/* Mocks for register handler */
func registerOk(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockUser, nil)
	authMock.EXPECT().EncodeTokenUserID(gomock.Any()).Return("tokentoken", nil)
	return authMock
}
func registerErrServer(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.User{}, errors.New("server error"))
	return authMock
}
func registerErrExist(ctrl *gomock.Controller) *mk.MockAuthenticator {
	authMock := mk.NewMockAuthenticator(ctrl)
	authMock.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.User{}, auth.ErrorUserAlreadyExist)
	return authMock
}
