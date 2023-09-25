package main

import (
	"time"
)

// Transfer related Types and Constructor
//
// -----------------------------------------------------------------------------
type CreateTransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

//-----------------------------------------------------------------------------

//	Account related Types and Constructor
//
//-----------------------------------------------------------------------------

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now().UTC(),
	}
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//-----------------------------------------------------------------------------

//	UserAccount related Types and Constructor
//
// -----------------------------------------------------------------------------
type UserAccount struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	PassHash string `json:"passHash"`
	SaltHash string `json:"saltHash"`
}

func NewUserAccount(email, passhash, salthash string) *UserAccount {
	return &UserAccount{
		Email:    email,
		PassHash: passhash,
		SaltHash: salthash,
	}
}

type CreateUserAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// -----------------------------------------------------------------------------
