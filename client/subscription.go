package client

import (
	"errors"
	"fmt"
	"net/textproto"
)

var errorInvalidOperation = errors.New("Invalid operation type")

// Subscription represents an abstraction
// over a TCP connection to consume a stream
// of events
type Subscription interface {
	// Read returns a message from the connection
	Read() ([]byte, error)

	// Close closes the underlying connection
	Close() error
}

type subscription struct {
	conn *textproto.Conn
}

// NewSubscription returns an instance of a Subscription
func NewSubscription(conn *textproto.Conn) Subscription {
	return &subscription{conn}
}

// Read reads from the underlying connection
// It does so by consuming the stream till a delimiter
// If the first line consumed is a
//   - MSG: The next line is read to return the message
//	 - ACK: Carries on with reading the stream
//	 - ERR: The next line is read and an error message is returned
//
// The connection is closed and an error is returned if reading
// from the connection fails at any point or if the command is not supported
//
// Note that this method blocks till it can read from the connection
func (s *subscription) Read() ([]byte, error) {
	for {
		line, err := s.conn.ReadLine()
		if err != nil {
			s.conn.Close()
			return nil, err
		}

		switch line {
		case MSG:
			data, err := s.conn.ReadLine()
			if err != nil {
				s.conn.Close()
				return nil, err
			}

			return []byte(data), nil
		case ACK:
			continue
		case ERR:
			errDesc, err := s.conn.ReadLine()
			if err != nil {
				s.conn.Close()
				return nil, err
			}

			return nil, fmt.Errorf("Error from server: %s", errDesc)
		default:
			s.conn.Close()
			return nil, errorInvalidOperation
		}
	}
}

// Close closes the underlying connection
// If there was an error encountered while closing the connection
// it is returned
func (s *subscription) Close() error {
	return s.conn.Close()
}
