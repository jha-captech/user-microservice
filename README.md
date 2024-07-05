# user-microservice Alternate Structures

## Introduction
This project is designed to demonstrate alternative structures for a basic microservice project. 

## Structures

### [Original](./original)
This is the original project and it used `cmd` for entry points and `internal` for internal packages. It also has both a lambda and a chi API entrypoint.

#### Pros
- stuff
#### Cons
- stuff

File layout:
```text
original
├── cmd
│   ├── http
│   │   ├── route
│   │   │   ├── health_check.go
│   │   │   ├── not_found.go
│   │   │   ├── route.go
│   │   │   ├── route_test.go
│   │   │   ├── user.go
│   │   │   └── utilities.go
│   │   ├── config.go
│   │   ├── http_requests.http
│   │   ├── logger.go
│   │   └── main.go
│   └── lambda
│       ├── handler
│       │   ├── encode_decode.go
│       │   ├── handler.go
│       │   └── user.go
│       ├── config.go
│       ├── lambda_requests.http
│       ├── logger.go
│       └── main.go
├── internal
│   ├── database
│   │   ├── entity
│   │   │   └── user.go
│   │   ├── database.go
│   │   ├── database_test.go
│   │   └── user.go
│   ├── testutil
│   │   ├── user.go
│   │   └── users_test.go
│   └── user
│       ├── mock
│       │   └── databasesession.go
│       ├── user.go
│       └── user_test.go
├── README.MD
├── docker-compose.yml
├── env.json
├── env.sample.json
├── go.mod
├── go.sum
├── http.dockerfile
├── makefile
├── samconfig.toml
└── template.yaml
```

TODOS:
- [ ] move all mocks out of test files and into `mock` packages.
- [ ] add lambda tests