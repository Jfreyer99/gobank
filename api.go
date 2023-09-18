package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// Using chi router seems to be more modern approach

	// Try organizing the code better by using subrouters to split those concerns apart
	//subrouter := router.PathPrefix("/").Subrouter()
	//getRoute := subrouter.HandleFunc("account", makeHTTPHandleFunc(s.handleAccount))

	// Try improving by using HandleFunc().Method("Get")

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods(http.MethodGet)

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleGetAccount)).Methods(http.MethodGet)

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleCreateAccount)).Methods(http.MethodPost)

	router.HandleFunc("/account/{id}", WithJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByID), s.store)).Methods(http.MethodGet)

	router.HandleFunc("/transfer",
		makeHTTPHandleFunc(s.handleTransfer)).Methods(http.MethodPost)

	router.HandleFunc("/account/{id}",
		makeHTTPHandleFunc(s.handleDeleteAccount)).Methods(http.MethodDelete)

	log.Println("JSON API Running on Port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	//account := NewAccount("jonas", "ff")
	id := GetID(r)

	account, err := s.store.GetAccountByID(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.GetAccounts()

	if err != nil {
		return WriteJSON(w, http.StatusNotFound, err)
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// Refactor and dont create JWT here only in Route handleCreateUserAccount
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := &CreateAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}

	defer r.Body.Close()

	account := NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName)
	err := s.store.CreateAccount(account)
	if err != nil {
		return err
	}

	tokenString, err := CreateJWT(account)
	if err != nil {
		return err
	}

	fmt.Println(tokenString)

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	if id < 1 {
		return fmt.Errorf("cannot Delete Account with ID less than 1")
	}

	rerr := s.store.DeleteAccount(id)
	if rerr != nil {
		return rerr
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferRequest := &CreateTransferRequest{}

	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}

	defer r.Body.Close()

	// ADD Logic to handle the transfer
	return WriteJSON(w, http.StatusOK, transferRequest)
}

//			JWT VALIDATION AND CREATION FOR NEW ACCOUNTS
//------------------------------------------------------------------------------------------

func WithJWTAuth(handlerFunc http.HandlerFunc, store Storage) http.HandlerFunc {
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

		account, err := store.GetAccountByID(claimID)
		if err != nil {
			PermissionDenied(w)
			return
		}

		id := GetID(r)
		if id != account.ID {
			PermissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func CreateJWT(account *Account) (string, error) {

	secret := os.Getenv("JWT_SECRET")

	mySigningKey := []byte(secret)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24).UTC()},
		IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
		Issuer:    "GoBank",
		ID:        strconv.Itoa(account.ID),
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

//-----------------------------------------------------------------------------------------------

//			Decorator for HandlerFunc and WriteJson Helper
//-----------------------------------------------------------------------------------------------

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusMethodNotAllowed, ApiError{Error: err.Error()})
		}
	}
}

//----------------------------------------------------------------------------------------------------------

// Extracted Functions
// ---------------------------------------------------------------------------------------------------------

func GetID(r *http.Request) int {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1
	}
	return id
}

func PermissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

// ---------------------------------------------------------------------------------------------------------
