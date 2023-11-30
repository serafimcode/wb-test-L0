package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	env "github.com/joho/godotenv"
	stanAdapter "github.com/serafimcode/wb-test-L0/stanAdapter"
)

func main() {
	if err := env.Load(); err != nil {
		log.Fatal(err)
	}
	clientId := os.Getenv("NATS_PUBLISHER_ID")
	connection, err := stanAdapter.InitStanConnection(clientId)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	orderBuffer := loadRawJson("publisher/model.json")
	channelName := os.Getenv("NATS_CHANNEL")
	err = connection.Publish(channelName, orderBuffer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message published")
}

func loadRawJson(path string) []byte {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	rawJson, err := os.ReadFile(absolutePath)
	if err != nil {
		log.Fatal(err)
	}
	return rawJson
}
