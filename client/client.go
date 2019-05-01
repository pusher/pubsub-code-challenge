package client

import (
	"errors"
	"fmt"
	"net/textproto"
)

// Response represents a response from the server
// If the server has sent an ACK, the Error is nil
// If the server has sent and ERR <desc> the Error
// property will contain the error description
type Response struct {
	Error error
}

// IsError indicates if the response was an error
func (r *Response) IsError() bool {
	return r.Error != nil
}

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
	Publish(channel string, data string) error

	// Subscribe to a channel
	Subscribe(channel string) (Subscription, error)

	// Raw can be used to send any arbitrary command to
	// the server and return a response
	Raw(cmd string) (*Response, error)
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
// An error is returned if there is either a connection error
// Errors returned by the server are returned as part of the response
func (c *client) Publish(channel string, data string) error {
	resp, err := c.Raw(fmt.Sprintf("%s\r\n%s\r\n%s", PUBLISH, channel, data))
	if err != nil {
		return err
	}

	if resp.IsError() {
		return resp.Error
	}

	return nil
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

	resp, err := readResponse(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if resp.IsError() {
		return nil, resp.Error
	}

	return NewSubscription(conn), nil
}

// Raw passes on any arbitrary command to the server
//
// Errors returned by this function are similar to that of
// Publish
func (c *client) Raw(cmd string) (*Response, error) {
	conn, err := textproto.Dial("tcp", c.address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = conn.PrintfLine("%s", cmd)
	if err != nil {
		return nil, fmt.Errorf("Failed to write to connection: %s", err.Error())
	}

	resp, err := readResponse(conn)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func readFromConnection(conn *textproto.Conn) (string, error) {
	data, err := conn.ReadLine()
	if err != nil {
		return "", fmt.Errorf("Failed to read from connection: %s", err.Error())
	}

	return data, nil
}

func readResponse(conn *textproto.Conn) (*Response, error) {
	resp, err := readFromConnection(conn)
	if err != nil {
		return nil, err
	}

	switch resp {
	case ACK:
		return &Response{}, nil
	case ERR:
		errDesc, err := readFromConnection(conn)
		if err != nil {
			return nil, err
		}
		return &Response{Error: errors.New(errDesc)}, nil
	default:
		return nil, errorInvalidOperation
	}
}
