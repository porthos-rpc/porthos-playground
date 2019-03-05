# Porthos Dashboard

A playground application for trying out a porthos-rpc cluster.

## Requirements

- RabbitMQ (any other AMQP 0.9.1 broker)
- SQLite

## Build

```sh
go build
```

## Running the app

After building the executable, you may run the playground as following:

```sh
./porthos-playground -bind 8080 -broker amqp:// -db playground.db
```

## Using the playground in your local docker environment

In the image, there's the `dockerize` utility so you can use it to start the playground only after RabbitMQ is up and running.

```
playground:
  image: porthos/porthos-playground
  command: dockerize -wait tcp://broker:5672 -timeout 60s /go/src/github.com/porthos-rpc/porthos-playground/playground
  links:
   - broker
  environment:
    BROKER_URL: amqp://guest:guest@broker:5672/
    BIND_ADDRESS: :8080
  ports:
   - "8080:8080"
```

## Local environment

To speed things up we provide a docker environment. Just run `docker-compose up` then open `http://localhost:8080/`.
