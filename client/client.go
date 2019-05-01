package client

import (
	"fmt"
	"net/textproto"
)

const (
	PUBLISH   = "PUB"
	SUBSCRIBE = "SUB"
	MSG       = "MSG"
	ACK       = "ACK"
	ERR       = "ERR"
)

// Client represents a publish/ subscribe client
type Client interface {
	// Publish a message to a channel
	Publish(channel string, data []byte) error

	// Subscribe to a channel
	Subscribe(channel string) (Subscription, error)

	// Raw can be used to send any arbitrary command to
	// the server
	Raw(cmd string) error
}

type client struct {
	address string
}

// Returns an instance of a Client
func New(address string) Client {
	return &client{address}
}

// Publish sends data to subscribers of a given channel
// by sending a `PUB` command to the underlying connection
// Note that a single connection is created to send the command
// and read a response after which the connection is closed
//
// An error is returned if
//   - There is either a connection error
//	 - The server fails to send an ACK
//	   and sends an ERR instead
func (c *client) Publish(channel string, data []byte) error {
	return c.Raw(fmt.Sprintf("%s\r\n%s\r\n%s", PUBLISH, channel, string(data)))
}

// Subscribe returns a Susbcription object that wraps the underlying
// connection and provides an abstraction to interact with it
// It continues to hold a reference to the connection so consumers
// may continue to read from it and consume events
//
// An error is returned if
//   - There is an error establishing the connection
// 	 - There is an error writing the SUBSCRIBE command
//	   to the connection
//   - There is an error reading a response from the
//	   connection
func (c *client) Subscribe(channel string) (Subscription, error) {
	conn, err := textproto.Dial("tcp", c.address)
	if err != nil {
		return nil, err
	}

	err = conn.PrintfLine("%s\r\n%s", SUBSCRIBE, channel)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Failed to write to connection: %s", err.Error())
	}

	err = readResponse(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return NewSubscription(conn), nil
}

// Raw passes on any arbitrary command to the server
//
// Errors returned by this function are similar to that of
// Publish
func (c *client) Raw(cmd string) error {
	conn, err := textproto.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.PrintfLine("%s", cmd)
	if err != nil {
		return fmt.Errorf("Failed to write to connection: %s", err.Error())
	}

	err = readResponse(conn)
	if err != nil {
		return err
	}

	return nil
}

func readResponse(conn *textproto.Conn) error {
	resp, err := conn.ReadLine()
	if err != nil {
		return fmt.Errorf(
			"Failed to read from connection: %s",
			err.Error(),
		)
	}

	switch resp {
	case ACK:
		return nil
	case ERR:
		errDesc, err := conn.ReadLine()
		if err != nil {
			return fmt.Errorf("Failed to read from connection: %s", err.Error())
		}
		return fmt.Errorf("Error from server: %s", errDesc)
	default:
		return errorInvalidOperation
	}
}
