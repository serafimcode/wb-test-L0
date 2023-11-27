package main

import (
	"fmt"
	env "github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := env.Load(); err != nil {
		log.Fatal("Error: loading .env file")
	}

	clusterId := os.Getenv("NATS_CLUSTER_ID")
	clientId := os.Getenv("NATS_SUBSCRIBER_ID")
	channelName := os.Getenv("NATS_CHANNEL")
	natsUrl := os.Getenv("NATS_URL")

	sc, err := stan.Connect(clusterId, clientId, stan.NatsURL(natsUrl))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	sub, err := sc.Subscribe(channelName, func(msg *stan.Msg) {
		fmt.Printf("Subscriber get message: %v", string(msg.Data))
	})

	if err != nil {
		log.Fatal(err)
	}
	defer sub.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	http.ListenAndServe(":8080", nil)
}
