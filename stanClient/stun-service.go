package stanClient

import (
	"encoding/json"
	"fmt"
	"github.com/serafimcode/wb-test-L0/dbAdapter"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"

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
	_, err = connection.Subscribe(channelName, ordersSubscribeHandler)

	if err != nil {
		log.Fatal(err)
	}
	return &stanService
}

func ordersSubscribeHandler(msg *stan.Msg) {
	fmt.Println("Subscriber received message")

	orderInfo, err := validateOrderInfoDTO(msg.Data)
	if err != nil {
		log.Println("JSON validation error:", err)
		return
	}

	orderId := orderInfo.OrderUID
	orderInfoRaw, err := json.Marshal(orderInfo)
	var orderInfoJson map[string]interface{}
	err = json.Unmarshal(orderInfoRaw, &orderInfoJson)
	if err != nil {
		log.Println("JSON unmarshal error:", err)
		return
	}

	order := model.Order{ID: orderId, Info: orderInfoJson}
	orderRepository := repository.OrderRepository{Db: dbAdapter.GetDb()}
	orderRepository.Create(order)
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
