package main

import (
	"fmt"
	env "github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
	"log"
	"os"
)

func main() {

	orderBuffer, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := env.Load("../.env"); err != nil {
		log.Fatal(err)
	}

	clusterId := os.Getenv("NATS_CLUSTER_ID")
	clientId := os.Getenv("NATS_PUBLISHER_ID")
	channelName := os.Getenv("NATS_CHANNEL")
	natsUrl := os.Getenv("NATS_URL")

	sc, err := stan.Connect(clusterId, clientId, stan.NatsURL(natsUrl))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	err = sc.Publish(channelName, orderBuffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Message published")
}
