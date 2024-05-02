# API GO

## Why

1. To learn go: it's my first project in golang.
2. To learn how to use TestContainers.
3. To experiment with the different deployment option available.
4. To create a real-scenario helm chart (with a db dependency).
5. To extract a starter from this chart.
6. To avoid doing real life task.

## Deploy development environment

The development environment is just a mysql docker container, using your `.env.dev` environment variable.

1. Create your `.env.dev` file:
```properties
MYSQL_DATABASE=
MYSQL_ROOT_PASSWORD=
MYSQL_USER=
MYSQL_PASSWORD=
MYSQL_HOST=
MYSQL_PORT=
```

2. Run the the docker compose dev command `docker compose -f docker-compose.dev.yaml up`.

3. Run the go project `go run api.go`.

The project should be running on `127.0.0.1:8080`.

## Deploy to production

Create your `.env.prod` file:
```properties
MYSQL_DATABASE=
MYSQL_ROOT_PASSWORD=
MYSQL_USER=
MYSQL_PASSWORD=
MYSQL_HOST=
MYSQL_PORT=
```

### Docker compose

run `docker compose -f docker-compose.prod.yaml up --build --force-recreate`
The project should be running on `127.0.0.1:8080`.

### Helm
INCOMMING

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
- [x] Containerized the application
- [x] Create a deployment docker-compose file
- [ ] Create a helm chart
