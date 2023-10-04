package api

import (
	"net/http"
)

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

	defer r.Body.Close()

	if err := s.store.DeleteUserAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, "{success: true}")
}
