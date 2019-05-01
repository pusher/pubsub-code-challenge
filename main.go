package main

import (
	"flag"

	"github.com/pusher/pubsub-coding-challenge-test-suite/client"
	"github.com/pusher/pubsub-coding-challenge-test-suite/test"

	"github.com/fatih/color"
)

func main() {
	address := flag.String("address", "localhost:8081", "Host running the pubsub TCP server")
	pubSubClient := client.New(*address)

	for testName, fn := range test.Tests {
		color.White("Running test -> (%s)", testName)
		err := fn(pubSubClient)
		if err != nil {
			color.Red("Test failed (%s)", err.Error())
		} else {
			color.Green("Test passed")
		}
	}
}
