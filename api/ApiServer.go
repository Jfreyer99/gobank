package api

import (
	"log"
	"net/http"

	"github.com/Jfreyer99/gobank/storage"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
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

	router.HandleFunc("/userAccount", makeHTTPHandleFunc(s.handleCreateUserAccount)).Methods(http.MethodPost)

	router.HandleFunc("/account/{id}/{number}", WithJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByIDAndNumber), s.store)).Methods(http.MethodGet)

	router.HandleFunc("/transfer",
		makeHTTPHandleFunc(s.handleTransfer)).Methods(http.MethodPost)

	router.HandleFunc("/account/{id}",
		makeHTTPHandleFunc(s.handleDeleteAccount)).Methods(http.MethodDelete)

	log.Println("JSON API Running on Port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}
