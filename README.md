# pubsub coding challenge test suite

This repository contains the test suite to run against interview solutions to the [pubsub-coding-challenge](https://github.com/pusher/pubsub-coding-challenge)

## Building

There are scripts located in the `scripts` directory that will build the test suite for a few different
architectures and operating systems. Executables will be output to `target/<arch>`.

## Test suite options

- `-target`: This points to the address of the server that is running the pubsub system. It defaults
to `localhost:8081`
