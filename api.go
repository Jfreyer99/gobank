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

	// Try organizing the code better by using subrouters to split those concerns apart
	//subrouter := router.PathPrefix("/").Subrouter()
	//getRoute := subrouter.HandleFunc("account", makeHTTPHandleFunc(s.handleAccount))

	// Try improving by using HandleFunc().Method("Get") 
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccount))

	log.Println("JSON API Running on Port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error{
	if r.Method == "GET" {
		return s.handleGetAccount(w,r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w,r)
	}	
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w,r)
	}

	return fmt.Errorf("method not allowed: %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error{
	//account := NewAccount("jonas", "ff")
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil{
		return err
	}

	//TODO getUser from store and pass it to WriteJson function
	account, err := s.store.GetAccountByID(id)

	if err != nil{
		return WriteJSON(w, http.StatusNotFound, err)
	}

	return WriteJSON(w, http.StatusOK, account)
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