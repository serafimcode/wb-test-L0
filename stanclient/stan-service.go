package stanclient

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"

	"github.com/serafimcode/wb-test-L0/dbadapter"
	"github.com/serafimcode/wb-test-L0/model"
	"github.com/serafimcode/wb-test-L0/repository"
)

var stanService = StanService{}

type StanService struct {
	Connection stan.Conn
	Subs       []stan.Subscription
}

func InitStan() *StanService {
	clientId := os.Getenv("NATS_SUBSCRIBER_ID")
	connection, err := InitStanConnection(clientId)
	if err != nil {
		log.Fatal(err)
	}
	stanService.Connection = connection

	channelName := os.Getenv("NATS_ORDERS_CHANNEL")
	_, err = connection.Subscribe(channelName, ordersSubscribeHandler, stan.DurableName("orders-durable"))

	if err != nil {
		log.Fatal(err)
	}
	return &stanService
}

func ordersSubscribeHandler(msg *stan.Msg) {
	orderInfo, err := validateOrderInfoDTO(msg.Data)
	if err != nil {
		log.Println("JSON validation error:", err)
		return
	}

	orderId := orderInfo.OrderUID
	jsonInfo, err := json.Marshal(orderInfo)

	if err != nil {
		log.Println("JSON unmarshal error:", err)
		return
	}

	order := model.Order{ID: orderId, Info: jsonInfo}
	orderRepository := repository.OrderRepository{Db: dbadapter.GetDb()}
	orderRepository.Create(&order)
}

func validateOrderInfoDTO(data []byte) (*model.OrderInfoDTO, error) {
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
