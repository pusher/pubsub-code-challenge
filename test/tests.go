package test

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/pusher/pubsub-code-challenge/client"
)

// Maintains a map of the test name to a function that accepts the client and returns an error
// A test is successful if there is no error returned
var Tests = map[string]func(cli client.Client) error{
	"Subscribe to channel `foo`/ publish to channel `foo`": func(cli client.Client) error {
		channelName := "foo"

		sub, err := cli.Subscribe(channelName)
		if err != nil {
			return err
		}
		defer sub.Close()

		data := "Hello"
		err = cli.Publish(channelName, data)
		if err != nil {
			return err
		}

		msg, err := sub.Read()
		if err != nil {
			return err
		}

		if msg != data {
			return fmt.Errorf("Expected message to be `Hello`, but got %s", msg)
		}

		return nil
	},
	"Subscribe to channel `foo`/ publish to channel `bar`": func(cli client.Client) error {
		timeout := 1 * time.Second
		sub, err := cli.Subscribe("foo")
		if err != nil {
			return err
		}
		defer sub.Close()

		data := "Hello"
		err = cli.Publish("bar", data)
		if err != nil {
			return err
		}

		msgChan := make(chan string, 1)
		errChan := make(chan error, 1)
		go func() {
			msg, err := sub.Read()
			if err != nil {
				errChan <- err
				return
			}

			msgChan <- msg
		}()

		select {
		case <-msgChan:
			return errors.New("Expected to not recieve a message on channel `foo`, but got one")
		case err := <-errChan:
			return err
		case <-time.After(timeout):
			return nil
		}

		return nil
	},
	"Multiple subscriptions to channel `foo`/ publish to channel `foo`": func(cli client.Client) error {
		numSubscriptions := 3
		channelName := "foo"

		msgChan := make(chan string, numSubscriptions)
		errChan := make(chan error, 1)

		for i := 0; i < numSubscriptions; i++ {
			go func() {
				sub, err := cli.Subscribe(channelName)
				if err != nil {
					errChan <- err
					return
				}
				defer sub.Close()

				msg, err := sub.Read()
				if err != nil {
					errChan <- err
					return
				}

				msgChan <- msg
			}()
		}

		// grace to allow subscriptions to settle
		time.Sleep(500 * time.Millisecond)

		data := "Hello"
		err := cli.Publish(channelName, data)
		if err != nil {
			return err
		}

		msgCount := 0
		for {
			select {
			case msg := <-msgChan:
				if msg == data {
					msgCount += 1
				}

				if msgCount == numSubscriptions {
					return nil
				}
			case err := <-errChan:
				return err
			}
		}

		return fmt.Errorf("Expected 3 messages to have been recieved, but got %d", msgCount)
	},
	"Subscribe to channel `foo`/ multiple publishes to channel `foo`": func(cli client.Client) error {
		channelName := "foo"
		numPublishes := 3

		sub, err := cli.Subscribe(channelName)
		if err != nil {
			return err
		}
		defer sub.Close()

		data := "Hello"
		for i := 0; i < numPublishes; i++ {
			err := cli.Publish(channelName, data)
			if err != nil {
				return err
			}
		}

		successfulReads := 0
		for {
			msg, err := sub.Read()
			if err != nil {
				return err
			}

			if msg == data {
				successfulReads += 1
			}

			if successfulReads == numPublishes {
				break
			}
		}

		if successfulReads != numPublishes {
			return fmt.Errorf(
				"Expected subscription to `foo` to recieve 3 messages, but got %d",
				successfulReads,
			)
		}

		return nil
	},
	"Invalid input": func(cli client.Client) error {
		resp, err := cli.Raw("FOO\r\nchannel")
		if err != nil {
			return err
		}

		if !resp.IsError() {
			return fmt.Errorf("Expected an error from the server, but got none")
		}

		return nil
	},
	"Subscribe to channel `foo`/ high rate of publishes to channel `foo`": func(cli client.Client) error {
		channelName := "foo"
		numPublishes := 100

		sub, err := cli.Subscribe(channelName)
		if err != nil {
			return err
		}
		defer sub.Close()

		msgChan := make(chan string, numPublishes)
		errChan := make(chan error, 1)
		go func() {
			for {
				msg, err := sub.Read()
				if err != nil {
					errChan <- err
					return
				}

				msgChan <- msg
			}
		}()

		data := "Hello"
		for i := 0; i < numPublishes; i++ {
			err := cli.Publish(channelName, data)
			if err != nil {
				return err
			}
		}

		msgCount := 0
		for {
			select {
			case err := <-errChan:
				return err
			case msg := <-msgChan:
				if msg == data {
					msgCount += 1
				}

				if msgCount == numPublishes {
					return nil
				}
			}
		}

		return fmt.Errorf("Expected %d messages to be read, but got %d", numPublishes, msgCount)
	},
	"Subscribe to channel `foo` at a high rate/ publish uniquely to each": func(cli client.Client) error {
		numSubscriptions := 200

		msgChan := make(chan string, numSubscriptions)
		subErrChan := make(chan error, 1)
		pubErrChan := make(chan error, 1)

		wg := sync.WaitGroup{}
		wg.Add(numSubscriptions)

		for i := 0; i < numSubscriptions; i++ {
			go func(idx int) {
				sub, err := cli.Subscribe("chan" + fmt.Sprint(idx))
				wg.Done()
				if err != nil {
					subErrChan <- err
					return
				}
				defer sub.Close()

				msg, err := sub.Read()
				if err != nil {
					subErrChan <- err
					return
				}

				msgChan <- msg
			}(i)
		}

		wg.Wait()

		for i := 0; i < numSubscriptions; i++ {
			go func(idx int) {
				err := cli.Publish("chan"+fmt.Sprint(idx), fmt.Sprint(idx))
				if err != nil {
					pubErrChan <- err
					return
				}
			}(i)
		}

		receivedMsgCount := 0
		for {
			select {
			case err := <-subErrChan:
				return err
			case err := <-pubErrChan:
				return err
			case <-msgChan:
				receivedMsgCount += 1

				if receivedMsgCount == numSubscriptions {
					return nil
				}
			}
		}

		return fmt.Errorf("Expected %d messages to be recieved, but got %d", numSubscriptions, receivedMsgCount)
	},
}
