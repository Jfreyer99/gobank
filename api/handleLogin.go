package api

import "net/http"

// TODO
// CREATE loginrequest struct
// CHECK Headers for x-jwt-token to authenticate
// IF not headers provided check database for user and compare password hashes
// IF they match create new jwt and send back to client to store as cookie
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	return nil
}
