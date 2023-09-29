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
	p, err := HashPassword(createUserAccountRequest.Password)

	userAccount := types.NewUserAccount(createUserAccountRequest.Email, p)
	err = s.store.CreateUserAccount(userAccount)
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

func (s *APIServer) handleGetUserAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := GetNumberParam(r, "id")
	if err != nil {
		return nil
	}

	account, err := s.store.GetUserAccountByID(id)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetUserAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetUserAccounts()
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleDeleteUserAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetNumberParam(r, "id")
	if err != nil {
		return nil
	}

	if err := s.store.DeleteUserAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, "{success: true}")
}
