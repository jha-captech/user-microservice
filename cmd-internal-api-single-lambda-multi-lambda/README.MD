# user-microservice

---

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

---

## App: Web API

### Dev

API running locally.

#### Start
```cmd
make http_dev
```

### Dev - Docker
API and DB running in docker.

#### Up
```cmd
make app_up
```

#### Down
```cmd
make app_down
```

---

## App: Lambda

### Run Single Lambda With Multiple Routes
```cmd  
make single_lambda_local_api
```

### Run Multiple Lambdas, With 1 Route per Lambda
```cmd  
make single_lambda_local_api
```
