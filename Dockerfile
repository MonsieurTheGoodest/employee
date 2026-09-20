FROM golang:tip-trixie

WORKDIR /employee

COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY ./config ./config
COPY go.mod go.sum ./

ENV INITDB_PATH=./internal/repository/initdb/schema.sql
ENV CONFIG_PATH=./config/config.yml

RUN go mod download

ENTRYPOINT ["go", "run", "./cmd/main.go"]