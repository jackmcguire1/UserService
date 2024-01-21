# User Service

[![Go Report Card](https://goreportcard.com/badge/github.com/jackmcguire1/UserService)](https://goreportcard.com/report/github.com/jackmcguire1/UserService)
[![codecov](https://codecov.io/gh/jackmcguire1/UserService/graph/badge.svg?token=URT8YBBJFF)](https://codecov.io/gh/jackmcguire1/UserService)

[git]:    https://git-scm.com/
[golang]: https://golang.org/
[modules]: https://github.com/golang/go/wiki/Modules
[goLand]: https://www.jetbrains.com/go/
[golint]: https://github.com/golangci/golangci-lint
[docker]: https://www.docker.com/products/docker-desktop

## ABOUT
> This repo contains a go module that exposes a User Microservice using MongoDB as a datastore

### Prerequisites

- [Git][git]
- [Go 1.21.1][golang]+
- [Docker][docker]

### SWAGGER HUB
> https://app.swaggerhub.com/apis/jackmcguire1/User-Service/1.0.2

### SETUP
> setup your mongo connection details in docker-compose.yaml
```yaml
environment:
    - MONGO_HOST=mongodb+srv://****
    - MONGO_DATABASE=***
    - MONGO_USERS_COLLECTION=users
  ```

#### run the docker-compose stack
```shell
docker-compose up -d && docker compose watch
```

### Environment Variables
- EVENTS_URL - external HTTP endpoint provided by interested services
- LOG_VERBOSITY - warn | error | info | debug
- MONGO_HOST - your mongo host url
- MONGO_DATABASE - your mongo database
- MONGO_USERS_COLLECTION - your mongo user's collection

## REQUIREMENTS
The service must allow you to:
- add a new User
- modify an existing User
- remove a User
- return a list of the Users, allowing for filtering by certain criteria (e.g. all Users with the
country &quot;GB&quot;)

The service must include:
- A sensible storage mechanism for the Users
- The ability to send events to notify other interested services of changes to User entities

## Thanks

This project exists thanks to **ME!**.

## Donations
All donations are appreciated!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](http://paypal.me/crazyjack12)
