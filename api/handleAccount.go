package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleGetAccountByIDAndNumber(w http.ResponseWriter, r *http.Request) error {
	//account := NewAccount("jonas", "ff")
	id := GetID(r)
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

// Add number to request Parameter for unique identification of Account PK(account_id, account_number)
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	numberStr := mux.Vars(r)["number"]
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return err
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
