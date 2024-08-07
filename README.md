# user-microservice Alternate Structures

## Introduction
This project is designed to demonstrate alternative structures for a basic microservice project. 

## Structures

### [Original](./original)
This is the original project and it used `cmd` for entry points and `internal` for internal packages. It also has both a lambda and a chi API entrypoint.

#### Pros
- Well organized. 
- `cmd` is a good way to facilitate multiple entrypoint for the app. 
- 
#### Cons
- It is more complex and can be hard to follow the flow of the application.
- combining both and API and a Lambda into the same project means that there is probably more code in `cmd` than there really should be.
- Due to the fact that the app is composed of multiple packages, there are abstractions between lairs which can be both a good and bad thing. 

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
│       │   └── database.go
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
- [ ] finish adding swagger support with `swaggo`

---

### [Flat](./flat)
This is a simplified version of the original project that uses a flat project structure. This version also uses all standard library packages, with the exception of the Postgres database driver

#### Pros
- Very fast to get up and running.
- Simple to understand where everything is. 
- No need to manage imports from other packages.
- Using all standard library packages is pretty powerful. 
  - Using `database/sql` instead of an ORM like `GORM` is actually pretty nice as you can see exactly what is going on with your SQL queries.
  - `net/http` is more than enough for most basic things. `net/http` handlers are also compatible with routing libraries like `go-chi/chi`.
#### Cons
- It can get a little crowded with all the files in the root of the project.
- Not practical for a large project or a project with more than a couple of devs. 
- Lack of abstractions may cause issues later on. 

File layout:
```text
flat
├── Dockerfile
├── README.MD
├── config.go
├── database.go
├── docker-compose.yml
├── encode_decode.go
├── go.mod
├── go.sum
├── handlers.go
├── main.go
├── makefile
├── middleware.go
├── models.go
├── postgres_setup.sql
├── requests.http
├── routes.go
└── users.go
```

TODOS:
- [ ] add tests

---

### [`cmd` and `internal` folders - API only](./cmd-internal-api-only)
This is a reorganized version of the flat project structure optimized for a multi person team. As such, it utilizes an entrypoint in `cmd` and logic is defined in packages inside of `internal`. 

#### Pros
- Structure facilitates multiple team members working on the project at once.
#### Cons
- Not as simple as the flat structure and requires packages to be imported from across the package.

File layout:
```text
cmd-internal-api-only
├── cmd
│   └── api
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── database.go
│   ├── middleware
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── middleware.go
│   │   └── recovery.go
│   ├── models
│   │   └── models.go
│   ├── server
│   │   ├── encode_decode.go
│   │   ├── handlers.go
│   │   ├── health_handlers.go
│   │   ├── routes.go
│   │   └── user_handlers.go
│   └── user
│       └── users.go
├── Dockerfile
├── README.MD
├── docker-compose.yml
├── go.mod
├── go.sum
├── makefile
├── postgres_setup.sql
└── requests.http
```

TODOS:
- [ ] add tests

---

### [`cmd` and `internal` folders - API and Lambda](./cmd-internal-api-and-lambda)
An alternate version of "`cmd` and `internal` folders - API only" that also adds a lambda that uses a revers proxy with existing API handlers.

#### Pros
- It has lambdas
#### Cons
- All endpoints are in a single lambda so there is some tighter coupling than might be desired. 

File layout:
```text
cmd-internal-api-and-lambda
├── cmd
│   ├── api
│   │   └── main.go
│   └── lambda
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── database.go
│   ├── middleware
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── middleware.go
│   │   └── recovery.go
│   ├── models
│   │   └── models.go
│   ├── server
│   │   ├── encode_decode.go
│   │   ├── handlers.go
│   │   ├── health_handlers.go
│   │   ├── routes.go
│   │   └── user_handlers.go
│   └── user
│       └── users.go
├── Dockerfile
├── README.MD
├── docker-compose.yml
├── env.json
├── env.sample.json
├── go.mod
├── go.sum
├── makefile
├── postgres_setup.sql
├── requests.api.http
├── requests.lambda.http
├── samconfig.toml
└── template.yaml
```

TODOS:
- [ ] add tests

---

### [`cmd` and `internal` folders - Lambda Only](./cmd-internal-lambda-only)
An alternate version of "`cmd` and `internal` folders - API only" that also only has lambdas and no API. Each lambda represents a single endpoint.

#### Pros
- It has lambdas and each lambda represents a single endpoint, allowing for each endpoint to be scaled interdependently of each other.
#### Cons
- Multiple lambdas meaning each endpoint is deployed seperated.

File layout:
```text
cmd-internal-lambda-only
├── cmd
│   └── lambda
│       ├── create
│       │   └── main.go
│       ├── delete
│       │   └── main.go
│       ├── fetch
│       │   └── main.go
│       ├── list
│       │   └── main.go
│       └── update
│           └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── database.go
│   ├── handler
│   │   ├── handler.go
│   │   ├── response.go
│   │   └── user.go
│   ├── models
│   │   └── models.go
│   └── user
│       └── users.go
├── Dockerfile
├── README.MD
├── docker-compose.yml
├── env.json
├── env.sample.json
├── go.mod
├── go.sum
├── makefile
├── postgres_setup.sql
├── requests.lambda.http
├── samconfig.toml
└── template.yamll
└── template.yaml
```

TODOS:
- [ ] add tests