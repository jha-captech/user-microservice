# user-microservice

## Database: PostgresSQL In Docker

### Setup And Seed Database
```cmd
make db_seed
```

### Start Database
```cmd
make up
```

### Stop Database
```cmd
make down
```

## App: Web API

### Dev
API running locally.
#### Start[README.md](..%2FREADME.md)
```cmd
make http_dev
```

### Dev - Docker
API and DB running in docker.

#### Up
```cmd
make app_up
```
[README.md](..%2FREADME.md)
#### Down
```cmd
make app_down
```
