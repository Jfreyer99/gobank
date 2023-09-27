package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jfreyer99/gobank/types"
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

	account, err := s.store.GetAccountByIDAndNumber(id, number)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.GetAccounts()

	if err != nil {
		return WriteJSON(w, http.StatusNotFound, err)
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// GET USER_ID FROM TOKEN AND USE THAT ID TO CREATE USER AND PROTECT THE ROUTE USING THE MIDDLEWARE
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accountRequest := &types.CreateAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(accountRequest); err != nil {
		return err
	}
	defer r.Body.Close()

	account := types.NewAccount(accountRequest.FirstName, accountRequest.LastName, accountRequest.ID)

	err := s.store.CreateAccount(account)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)

}

// Add number to request Parameter for unique identification of Account PK(account_id, account_number)
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := GetNumberParam(r, "id")
	if err != nil {
		return err
	}

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
