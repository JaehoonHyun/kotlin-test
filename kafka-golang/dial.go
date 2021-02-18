package main

import (
	"net"
	"strconv"

	kafka "github.com/segmentio/kafka-go"
)

func Dial(network, address string) (*kafka.Conn, *kafka.Conn, error) {
	//먼저 cluster에 다이얼을 건다.
	// conn, err := kafka.Dial("tcp", "kafka-headless.banzaic-kafka:29092")
	conn, err := kafka.Dial(network, address)
	if err != nil {
		return nil, nil, err
	}

	//그럼 거기서 어떤 broker와 통신해야하는지 또 다이얼을 날린다
	controller, err := conn.Controller()
	if err != nil {
		return nil, nil, err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return nil, nil, err
	}

	return conn, controllerConn, nil
}
