package main

import (
	"log"
	//"fmt"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil{
		log.Fatal(err)
	}

	if err := store.printAccountTable(); err != nil{
		log.Fatal(err)
	}

	//fmt.Printf("no error so far")

	server := NewAPIServer(":3000", store)
	server.Run()
}