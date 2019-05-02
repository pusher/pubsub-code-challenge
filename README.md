# pubsub coding challenge test suite

This repository contains the test suite to run against interview solutions to the [pubsub-coding-challenge](https://github.com/pusher/pubsub-coding-challenge)

## Installing dependencies

Make sure you have `dep` installed by following instructions from [here](https://github.com/golang/dep#installation).

Then run

```
dep ensure
```

## Building

There are scripts located in the `scripts` directory that will build the test suite for a few different
architectures and operating systems. Executables will be output to `target/<arch>`.

## Options

- `-address`: This points to the address of the server that is running the pubsub system. It defaults
to `localhost:8081`.

## Tests

Currently the following test cases are implemented:

- **Simple publish**: Subscribing to a channel `foo` followed by a publish to it.

  The test passes if the message published to channel `foo` can be read by a subsciber
  listening on channel `foo`.
- **Publishing to an unsubscribed channel**: Subscribing to a channel `foo`, but publishing to
  channel `bar`.

  The test passes if a subscriber listening on channel `bar` _does not_ receive a message
  published to channel `foo`.
- **Multiple subscriptions**: Creating multiple subscriptions to channel `foo` followed by a publish
  to it.

  The test passes if all subscriptions made to channel `foo` receive the same message published
  to it.
- **Multiple publishes**: Creating a single subscription to channel `foo` but publishing mutliple times.

  The test passes succeed if the same subscription to channel `foo` receives all publishes made to
  channel `foo`.
- **Invalid input**: Sending unsupported operations to the server.

  The test passes if the server returns an appropriate unsupported operation error.
- **High rate of publishes**: Subscribing to a channel `foo` but publishing rapidly to it.

  The test passes if the subscription to channel `foo` receives all publishes made to it.
- **High rate of subscriptions**: Subscribing to several channels rapidly and publishing unique
  messages to each of them.

  The test passes if each channel recieves the single message that was published to it.

**Note**: Tests will time out after 5 seconds if they do not complete.