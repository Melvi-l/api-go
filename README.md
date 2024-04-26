# API GO

## Why

1. To learn go: it's my first project in golang.
2. To learn how to use TestContainers.
3. To experiment with the different deployment option available.
4. To create a real-scenario helm chart (with a db dependency).
5. To extract a starter from this chart.
6. To avoid doing real life task.

## Run development environment

- docker-compose: run the `docker compose up` command.

## Deploy

- docker-compose: INCOMMING
- helm: INCOMMING

## Env variable
```properties
MYSQL_DATABASE=
MYSQL_ROOT_PASSWORD=
MYSQL_USER=
MYSQL_PASSWORD=
MYSQL_HOST=
MYSQL_PORT=
```

# TODO
- [x] Add env variable dynamic configuration
- [x] Add test with [TestContainer](https://golang.testcontainers.org/)
- [x] Add a github action build and test workflows
