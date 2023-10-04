package api

import (
	"encoding/json"
	"net/http"

	"github.com/Jfreyer99/gobank/types"
)

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	createUserAccountRequest := &types.CreateUserAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(createUserAccountRequest); err != nil {
		return err
	}

	defer r.Body.Close()

	p, err := HashPassword(createUserAccountRequest.Password)

	userAccount := types.NewUserAccount(createUserAccountRequest.Email, p)
	err = s.store.CreateUserAccount(userAccount)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]bool{"success": false})
	}

	return WriteJSON(w, http.StatusOK, userAccount)
}
