package main

import (
	"crypto/rand"
	"log"
	"math/big"
	//"fmt"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// if err := store.DropTableAccount(); err != nil{
	// 	log.Fatal(err)
	// }

	defer store.Close()

	server := NewAPIServer(":3000", store)
	server.Run()
}

// Function used for generating JWT_SECRET env variable n >= 50
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-&%$"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret = append(ret, letters[num.Int64()])
	}
	return string(ret), nil
}
