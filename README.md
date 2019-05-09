# Chatkit team coding challenge

Welcome to the Pusher Chatkit team pre-interview challenge!

This challenge is preparation for the on-site interview. During the on-site
interview we will have a session where we discuss the problem and your
solution with you.

## How your solution will be evaluated

Your solution will be evaluated based on

- Correctness. The test suite in this repo will help guide you towards a
  correct solution, but it is by no means exhaustive.

- Code clarity and structure. We are very interested in how you write your
  code to make it clear to read, maintainable and easy to work on in future.

## How to submit your solution

We ask that you don't make your solution public (e.g. on Github).

Please submit a tarball or zipfile via email to your recruitment contact
inside Pusher.

## Requirements for submissions

In order that we can run your solution to test it, please work within the
following restrictions:

- If using Java, your code *must* run on OpenJDK
- If using C#, your code *must* run on Mono or .net Core (i.e. we must be able
  to run it on Ubuntu / MacOS)

To help us run your submission, please include:

- a bash script named `run` which will both compile and run your server,
  including fetching any required dependencies.
- a README file with your name and email, and any thoughts you have about the
  challenge

# The actual challenge

[Pub/sub](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern) is
the pattern on which Pusher's original (Channels) product built its success,
and one that remains core to many of the problems we solve today.

Because of this, we would like you to implement a very simple pub/sub bus.

A basic test suite is provided to help you evaluate what you have written. You
can find it [here](https://github.com/pusher/pubsub-code-challenge).

We expect the challenge to take approximately 3 hours to complete. Please keep
a rough eye on the time. It's perfectly fine to submit with known issues,
especially if you document them.

We would prefer you to complete the challenge in Go, Java or C#, though we
will accept submissions in other languages, drop us a line before you start so
that we can check that we'll have the expertise in house to evaluate your
work.

Note that we have excluded the major interpreted languages from the above set,
because they tend to raise different concerns than those you would encounter
working on our codebases.

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
suite of tests. You can find the list of tests that have been implemented
[here](./TESTS.md).

You can download a pre-built version of the test suite from the following
links:

- [Linux](https://github.com/pusher/pubsub-code-challenge/raw/master/dist/linux-amd64/pusher-interview-test)
- [MacOS](https://github.com/pusher/pubsub-code-challenge/raw/master/dist/darwin-amd64/pusher-interview-test)

Alternatively, if you already have the Go toolchain installed,

```
go get github.com/pusher/pubsub-code-challenge/pusher-interview-test
```

The source code of the test suite is available in this repo - [here, if you are
not reading the repo readme right
now](https://github.com/pusher/pubsub-code-challenge), and you are free to
examine it, and build the binary yourself if you wish. You will need a Go
toolchain to do so.

### Running

```
pusher-interview-test [-address <host:port>]
```
