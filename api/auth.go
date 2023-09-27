package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Jfreyer99/gobank/storage"
	"github.com/Jfreyer99/gobank/types"
	"github.com/golang-jwt/jwt/v5"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling JWT Auth Middleware")

		tokenString := r.Header.Get("x-jwt-token")
		token, err := ValidateJWT(tokenString)

		if err != nil {
			PermissionDenied(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			PermissionDenied(w)
			return
		}

		claimIDStr := claims["jti"].(string)
		claimID, err := strconv.Atoi(claimIDStr)
		if err != nil {
			PermissionDenied(w)
			return
		}

		account, err := store.GetUserAccountByID(claimID)
		if err != nil {
			PermissionDenied(w)
			return
		}

		id, err := GetNumberParam(r, "id")
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, fmt.Errorf("conversion went wrong"))
			return
		}
		if id != account.ID {
			PermissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func CreateJWT(userAccount *types.UserAccount) (string, error) {

	secret := os.Getenv("JWT_SECRET")

	mySigningKey := []byte(secret)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24).UTC()},
		IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
		Issuer:    "GoBank",
		ID:        strconv.Itoa(userAccount.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {

	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
