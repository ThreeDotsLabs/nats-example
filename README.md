# nats-example

This is an example Go project for running local development environment. It features:

* Hot code reloading with [reflex](https://github.com/cespare/reflex)
* Running multiple Docker containers with [Docker Compose](https://docs.docker.com/compose)
* Using [Go Modules](https://github.com/golang/go/wiki/Modules) for managing dependencies
* [Publisher](publisher/main.go) and [Subscriber](subscriber/main.go) services
* Communication over [NATS](https://nats.io/) between two services

Read more in our blog post: [Go Docker dev environment with Go Modules and live code reloading](https://threedots.tech/post/go-docker-dev-environment-with-go-modules-and-live-code-reloading/)

[![asciicast](https://asciinema.org/a/kas6dYKpMzyubpCmW9aOjnBIu.svg)](https://asciinema.org/a/kas6dYKpMzyubpCmW9aOjnBIu)

## Running

To start the services, simply run:

```
docker-compose up
```

Open new terminal and send some messages:

```
$ curl localhost:5000 -d "this is my message"
Sent message: this is my message with ID 01D09P02SBW5D0QPWP14QQZJWH
```

You should see the messages in the `subscriber` service output:

```
subscriber_1  | [00] 2019/01/03 11:01:27 received message: 01D09P02SBW5D0QPWP14QQZJWH, payload: this is my message
```
