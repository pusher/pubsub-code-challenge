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
