package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jfreyer99/gobank/types"
)

// Refactor and dont create JWT here only in Route handleCreateUserAccount
func (s *APIServer) handleCreateUserAccount(w http.ResponseWriter, r *http.Request) error {
	createUserAccountRequest := &types.CreateUserAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(createUserAccountRequest); err != nil {
		return err
	}

	defer r.Body.Close()

	// TODO Hash the password and salt acordindly and generate a proper salt using bycypt
	salt := "fasdfasdfasfdasfd"

	userAccount := types.NewUserAccount(createUserAccountRequest.Email, createUserAccountRequest.Password, salt)
	err := s.store.CreateUserAccount(userAccount)
	if err != nil {
		return err
	}

	tokenString, err := CreateJWT(userAccount)
	if err != nil {
		return err
	}

	fmt.Println(tokenString)

	return WriteJSON(w, http.StatusOK, userAccount)
}
