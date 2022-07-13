package middleware

import (
	"github.com/atrush/diploma.git/api/model"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"testing"
)

//  TestHandler_Login tests user register handler
func TestMiddlewareAuth(t *testing.T) {
	userID := uuid.New()

	//  jwtauth init
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	//  generate jwt token
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": userID})
	require.NoError(t, err)

	//  middleware reads Authorization header or jwt cookie to context
	jwtMiddleware := jwtauth.Verifier(tokenAuth)
	//  handler writes user_id to body
	toBodyHandler := writeUserIDToBody{t: t}

	tests := []TestMiddleware{
		{
			name:           "authenticated ok POST",
			middlewareFunc: []Middleware{MiddlewareAuth, jwtMiddleware},
			nextHandler:    toBodyHandler,
			method:         http.MethodPost,
			headers:        map[string]string{"Authorization": "Bearer " + tokenString},

			expectedBody: userID.String(),
			expectedCode: 200,
		},
		{
			name:           "authenticated ok GET",
			middlewareFunc: []Middleware{MiddlewareAuth, jwtMiddleware},
			nextHandler:    toBodyHandler,
			method:         http.MethodGet,
			headers:        map[string]string{"Authorization": "Bearer " + tokenString},

			expectedBody: userID.String(),
			expectedCode: 200,
		},
		{
			name:           "no jwt middleware - 401",
			middlewareFunc: []Middleware{MiddlewareAuth},
			nextHandler:    toBodyHandler,
			method:         http.MethodPost,
			headers:        map[string]string{"Authorization": "Bearer " + tokenString},

			expectedCode: 401,
		},
		{
			name:           "no Authorization header - 401",
			middlewareFunc: []Middleware{MiddlewareAuth, jwtMiddleware},
			nextHandler:    toBodyHandler,
			method:         http.MethodPost,

			expectedCode: 401,
		},
		{
			name:           "wrong token - 401",
			middlewareFunc: []Middleware{MiddlewareAuth, jwtMiddleware},
			nextHandler:    toBodyHandler,
			method:         http.MethodPost,
			headers:        map[string]string{"Authorization": "Bearer tokentoken"},

			expectedCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.CheckTest(t)
		})
	}
}

// writeUserIDToBody handler with testing.T writes user_id from context to body
type writeUserIDToBody struct {
	t *testing.T
}

func (wr writeUserIDToBody) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctxID := r.Context().Value(model.ContextKeyUserID).(string)
	require.NotEmpty(wr.t, ctxID)

	userID, err := uuid.Parse(ctxID)
	require.NoError(wr.t, err)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(userID.String())); err != nil {
		log.Fatal(err.Error())
	}
}
