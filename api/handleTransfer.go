package api

import (
	"encoding/json"
	"net/http"

	"github.com/Jfreyer99/gobank/types"
)

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferRequest := &types.CreateTransferRequest{}

	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}

	defer r.Body.Close()

	// ADD Logic to handle the transfer
	return WriteJSON(w, http.StatusOK, transferRequest)
}
