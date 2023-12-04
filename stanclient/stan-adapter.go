package stanclient

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"os"
)

func InitStanConnection(clientId string) (stan.Conn, error) {
	clusterId := os.Getenv("NATS_CLUSTER_ID")
	natsUrl := os.Getenv("NATS_URL")

	connection, err := stan.Connect(clusterId, clientId, stan.NatsURL(natsUrl))
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v succesfully connected to %v on %v\n", clientId, clusterId, natsUrl)
	return connection, nil
}
