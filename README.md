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
- **Publishing to an unsubscribed channel**: Subscribing to a channel `foo`, but publishing to
  channel `bar`.
- **Multiple subscriptions**: Creating multiple subscriptions to channel `foo` followed by a publish
  to it.
- **Multiple publishes**: Creating a single subscription to channel `foo` but publishing mutliple times.
- **Invalid input**: Sending unsupported operations to the server.
- **High rate of publishes**: Subscribing to a channel `foo` but publishing rapidly to it.

**Note**: Tests will timeout after 5 seconds if they do not complete.