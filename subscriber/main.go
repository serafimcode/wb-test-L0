package main

import (
	"encoding/json"
	"fmt"
	env "github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
	"github.com/serafimcode/wb-test-L0/model"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"

	"github.com/serafimcode/wb-test-L0/database"
	stanAdapter "github.com/serafimcode/wb-test-L0/stan"
)

var db *gorm.DB

func main() {
	if err := env.Load(); err != nil {
		log.Fatal(err)
	}

	db = database.InitDb()
	fmt.Println(db)

	clientId := os.Getenv("NATS_SUBSCRIBER_ID")
	connection, err := stanAdapter.InitStanConnection(clientId)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	channelName := os.Getenv("NATS_CHANNEL")
	sub, err := connection.Subscribe(channelName, func(msg *stan.Msg) {
		var orderInfo map[string]interface{}
		if err := json.Unmarshal(msg.Data, &orderInfo); err != nil {
			log.Fatal(err)
		}
		id := orderInfo["order_uid"].(string)
		record := model.Order{ID: id, Info: orderInfo}
		db.Create(&record)

		fmt.Printf("Subscriber get message: %v", orderInfo)
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
