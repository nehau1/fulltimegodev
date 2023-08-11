package main

import (
	"log"

	"github.com/Stiffjobs/toll-calculator/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://localhost:3000"
)

//Transport (HTTP, Kafka, gRPC, etc.) -> attach business logic to this transport

func main() {
	var (
		svc CalculatorServicer
		err error
	)

	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewHTTPClient(aggregatorEndpoint))

	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()
}
