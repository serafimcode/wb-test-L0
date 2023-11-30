package main

import (
	"fmt"
	"log"
	"net/http"

	env "github.com/joho/godotenv"

	"github.com/serafimcode/wb-test-L0/dbAdapter"
	"github.com/serafimcode/wb-test-L0/stanClient"
)

func main() {
	if err := env.Load(); err != nil {
		log.Fatal(err)
	}

	dbAdapter.InitDb()

	stanService := stanClient.InitStan()
	defer stanService.Connection.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})
	http.ListenAndServe(":8080", nil)
}
