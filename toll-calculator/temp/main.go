package main

import (
	"context"
	"log"
	"time"

	"github.com/Stiffjobs/toll-calculator/aggregator/client"
	"github.com/Stiffjobs/toll-calculator/types"
)

func main() {
	c, err := client.NewGRPCClient(":3001")
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 30.20,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
}
