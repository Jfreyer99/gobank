package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jfreyer99/gobank/types"
	"github.com/golang-jwt/jwt/v5"
)

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {

	loginReq := &types.LoginUserAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(loginReq); err != nil {
		return err
	}

	login := types.NewUserAccount(loginReq.Email, loginReq.Password)

	tokenString := r.Header.Get("x-jwt-token")

	if len(tokenString) == 0 {

		fmt.Println("no jwt token provided")

		userAccount, err := s.store.GetUserAccountByEmail(login.Email)
		if err != nil {
			return err
		}

		ok := CheckPasswordHash(userAccount.PassHash, login.PassHash)
		if !ok {
			fmt.Println("No user with password")
			return WriteJSON(w, http.StatusBadRequest, map[string]bool{"success": false})
		}

		tokenString, err := CreateJWT(userAccount)
		if err != nil {
			return WriteJSON(w, http.StatusBadRequest, map[string]bool{"success": false})
		}

		fmt.Println(tokenString)

		return WriteJSON(w, http.StatusOK, tokenString)
	}

	token, err := ValidateJWT(tokenString)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return WriteJSON(w, http.StatusBadRequest, map[string]bool{"success": false})
	}

	claimIDStr := claims["jti"].(string)
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]bool{"success": false})
	}

	acc, err := s.store.GetUserAccountByID(claimID)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]bool{"success": false})
	}

	return WriteJSON(w, http.StatusAccepted, acc)   
}