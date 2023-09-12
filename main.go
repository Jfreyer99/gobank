package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil{
		log.Fatal(err)
	}

	account := &Account{
		FirstName: "Vorname",
		LastName: "Nachname",
		Balance: 1000.56,
	}
	
	if err := store.CreateAccount(account); err != nil{
		log.Fatal(err)
	}

	if err := store.printAccountTable(); err != nil{
		log.Fatal(err)
	}

	//fmt.Printf("no error so far")

	server := NewAPIServer(":3000", store)
	server.Run()
}