# Porthos Dashboard

A playground application for trying out a porthos-rpc cluster.

## Requirements

- RabbitMQ (any other AMQP 0.9.1 broker)
- SQLite

## Build

First of all, you need to download `govendor`:

```sh
go get -u github.com/kardianos/govendor
```

Fetch all go dependencies:

```sh
govendor sync
```

Then build the app:

```sh
go build
```

## Running the app

After building the executable, you may run the playground as following:

```sh
./porthos-playground -bind 8080 -broker amqp:// -db playground.db
```

## Local environment

To speed things up we provide a docker environment. Just run `docker-compose up` then open `http://localhost:8080/`.
