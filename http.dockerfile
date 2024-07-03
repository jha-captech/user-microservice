FROM golang:1.22.3-alpine3.19 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./app ./cmd/http

FROM alpine:3.19 AS publish

WORKDIR /app

COPY --from=build ./app .

EXPOSE 8080

ENTRYPOINT [ "./app" ]