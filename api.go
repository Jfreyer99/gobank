package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer{
	return &APIServer{
		listenAddr: listenAddr,
		store: store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// Using chi router seems to be more modern approach
	
	// Try organizing the code better by using subrouters to split those concerns apart
	//subrouter := router.PathPrefix("/").Subrouter()
	//getRoute := subrouter.HandleFunc("account", makeHTTPHandleFunc(s.handleAccount))

	// Try improving by using HandleFunc().Method("Get") 
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleGetAccount)).Methods(http.MethodGet)

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleCreateAccount)).Methods(http.MethodPost)

	router.HandleFunc("/account/{id}", 
	makeHTTPHandleFunc(s.handleGetAccountByID)).Methods(http.MethodGet)

	router.HandleFunc("/account/{id}", 
	makeHTTPHandleFunc(s.handleDeleteAccount)).Methods(http.MethodDelete)

	log.Println("JSON API Running on Port: ", s.listenAddr)
	
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error{
	//account := NewAccount("jonas", "ff")
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil{
		return err
	}

	account, err := s.store.GetAccountByID(id)

	if err != nil{
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error{

	accounts, err := s.store.GetAccounts()

	if err != nil{
		return WriteJSON(w, http.StatusNotFound, err)
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error{
	createAccountRequest := &CreateAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil{
		return err
	}

	account := NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName)
	err:= s.store.CreateAccount(account)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}


func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error{

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil{
		return err
	}

	if id < 1{
		return fmt.Errorf("Cannot Delete Account with ID less than 1")
	}

	rerr := s.store.DeleteAccount(id)

	if rerr != nil{
	 	return rerr
	}

	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error{
	return nil
}


func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusMethodNotAllowed, ApiError{Error: err.Error()})
		}
	}
}