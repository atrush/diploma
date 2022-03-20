package api

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"log"
)

//  genJWTAuthToken generates jwtauth token from user id
func genJWTAuthToken(userID uuid.UUID) string {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": userID.String()})

	if err != nil {
		log.Fatal(err.Error())
	}
	return tokenString
}
