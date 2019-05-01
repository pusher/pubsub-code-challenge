package main

import (
	"flag"
	"time"

	"github.com/pusher/pubsub-code-challenge/client"
	"github.com/pusher/pubsub-code-challenge/test"

	"github.com/fatih/color"
)

const testTimeout = 5 * time.Second

func main() {
	address := flag.String("address", "localhost:8081", "Host running the pubsub TCP server")
	flag.Parse()

	pubSubClient := client.New(*address)

	for testName, fn := range test.Tests {
		errChan := make(chan error, len(test.Tests))

		color.White("Running test -> (%s)", testName)
		go func() {
			errChan <- fn(pubSubClient)
		}()

		select {
		case err := <-errChan:
			if err != nil {
				color.Red("Test failed (%s)", err.Error())
			} else {
				color.Green("Test passed")
			}
		case <-time.After(testTimeout):
			color.Red("Test timed out after waiting for %s", testTimeout.String())
		}
	}
}
