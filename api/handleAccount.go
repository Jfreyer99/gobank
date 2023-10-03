package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jfreyer99/gobank/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func (s *APIServer) handleGetAccountByIDAndNumber(w http.ResponseWriter, r *http.Request) error {
	//account := NewAccount("jonas", "ff")
	id, err := GetNumberParam(r, "id")
	if err != nil {
		return err
	}
	numberStr := mux.Vars(r)["number"]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return fmt.Errorf("no account number provided")
	}

	defer r.Body.Close()

	account, err := s.store.GetAccountByIDAndNumber(id, number)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.GetAccounts()

	defer r.Body.Close()

	if err != nil {
		return WriteJSON(w, http.StatusNotFound, err)
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// GET USER_ID FROM TOKEN AND USE THAT ID TO CREATE USER AND PROTECT THE ROUTE USING THE MIDDLEWARE
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accountRequest := &types.CreateAccountRequest{}

	tokenString := r.Header.Get("x-jwt-token")
	token, err := ValidateJWT(tokenString)

	if err != nil {
		return WriteJSON(w, http.StatusCreated, err)
	}

	defer r.Body.Close()

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return WriteJSON(w, http.StatusCreated, "{Not allowed}")
	}

	claimIDStr := claims["jti"].(string)
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		return WriteJSON(w, http.StatusCreated, err)
	}

	if err := json.NewDecoder(r.Body).Decode(accountRequest); err != nil {
		return err
	}

	account := types.NewAccount(accountRequest.FirstName, accountRequest.LastName, claimID)

	err = s.store.CreateAccount(account)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)
}

func (s *APIServer) handleGetAllAccount(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := GetNumberParam(r, "id")
	if err != nil {
		return err
	}

	defer r.Body.Close()

	number, err := GetNumberParam(r, "number")
	if err != nil {
		return nil
	}

	if id < 1 {
		return fmt.Errorf("cannot Delete Account with ID less than 1")
	}

	rerr := s.store.DeleteAccount(id, number)
	if rerr != nil {
		return rerr
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}
