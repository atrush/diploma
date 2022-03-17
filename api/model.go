package api

import "fmt"

type (
	LoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	contextKey string
)

func (r LoginRequest) Validate() error {
	if len(r.Login) < 3 {
		return fmt.Errorf("login must be larger then 3 symbols")
	}
	if len(r.Login) > 30 {
		return fmt.Errorf("login must be less then 30 symbols")
	}
	if len(r.Password) < 3 {
		return fmt.Errorf("password must be larger then 3 symbols")
	}
	return nil
}
