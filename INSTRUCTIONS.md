# Chatkit team coding challenge

[Pub/sub](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern) is
the pattern on which Pusher's original (Channels) product built its success,
and one that remains core to many of the problems we solve today.

Because of this, we would like you to implement a very simple pub/sub bus.

A basic test suite is provided to help you evaluate what you have written. You
can find it [here](TODO).

## Problem statement

You should write a program which listens on a TCP port for incoming
connections and accepts commands from clients on those connections. The two
commands are "pub(lish)" and "sub(scribe)". A more detailed specification of
the request format is described under "Protocol" below.

Both request types have a "channel name" as an argument.
Additionally the publish request has a payload as a second argument.

Whenever a publish request is issued, the payload included should be forwarded
to all connections who have previously subscribed to the same channel name.

Example:

- Client 1 connects
- Client 1 publishes "hello" to channel "world"
- *Nothing further happens, because there are no clients subscribed to
  channel "world"*
- Client 2 connects
- Client 2 subscribes to channel "world"
- Client 3 connects
- Client 3 publishes "anyone out there?" to channel "world"
- Client 2 receives "anyone out there?" from the server
- Client 4 connects
- Client 4 subscribes to channel "world"
- Client 3 publishes "can anybody hear me?" to channel "world"
- Clients 2 and 4 *both* receive "can anyone here me?" from the server

## Protocol

The network protocol is line oriented and takes place over TCP.

The server should listen on TCP port 8081 for connections.

On accepting a connection, the server should not send anything to the client

On connecting, the client should send a request to the server.

### Request format

Requests are made up of a command and some arguments.

The command and each argument are separated from each other with `\r\n`.

### Responses

The server responds to each command with an acknowledgement or error.

```
ACK
```

if the command is successful, or

```
ERR
description (argument)
```

if there was a problem with the request.

### Publish command

```
PUB
channel name (argument)
payload data (argument)
```

The server should respond with an ack or error, and then be ready to accept
another command.

### Subscribe command

```
SUB
channel name (argument)
```

The server should respond with an ack or error.

If the server responds with an ack, then no more commands may be issued on the
connection, and the server should start forwarding payloads published to the
specified channel like so:

```
MSG
payload data (argument)
```

If the server responded with an error, it should be ready to accept another
command.

## The test suite

### Installation

To evaulate your progress, or check your implementation, we've provided a
suite of tests.

You can download a pre-built version of the test suite from the following
links:

- [Linux](TODO)
- [MacOS](TODO)
- [Windows](TODO)

Alternatively, if you already have the Go toolchain installed,

```
go get github.com/pusher/pubsub-coding-challenge-test-suite
```

The source code of the test suite is available on
[Github](https://github.com/pusher/pubsub-coding-challenge-test-suite), and
you are free to examine it, and build the binary yourself if you wish. You
will need a Go toolchain to do so.

### Running


