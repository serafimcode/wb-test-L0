package main

import (
	"encoding/json"
	"fmt"
	"github.com/serafimcode/wb-test-L0/repository"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	env "github.com/joho/godotenv"
	"github.com/nats-io/stan.go"

	"github.com/serafimcode/wb-test-L0/dbAdapter"
	"github.com/serafimcode/wb-test-L0/model"
	"github.com/serafimcode/wb-test-L0/stanAdapter"
)

func main() {
	if err := env.Load(); err != nil {
		log.Fatal(err)
	}

	dbAdapter.InitDb()

	clientId := os.Getenv("NATS_SUBSCRIBER_ID")
	connection, err := stanAdapter.InitStanConnection(clientId)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	channelName := os.Getenv("NATS_CHANNEL")
	sub, err := connection.Subscribe(channelName, func(msg *stan.Msg) {
		orderInfo, err := validateStanData(msg.Data)
		if err != nil {
			log.Println("JSON validation error:", err)
			return
		}

		id := orderInfo.OrderUID
		orderJson, err := json.Marshal(orderInfo)
		record := model.Order{ID: id, Info: orderJson}

		orderRepository := repository.OrderRepository{Db: dbAdapter.GetDb()}
		orderRepository.Create(record)
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

func validateStanData(data []byte) (*model.OrderInfoDTO, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	var orderInfo model.OrderInfoDTO

	if err := json.Unmarshal(data, &orderInfo); err != nil {
		return nil, err
	}
	err := validate.Struct(orderInfo)
	if err != nil {
		return nil, err
	}
	return &orderInfo, nil
}
