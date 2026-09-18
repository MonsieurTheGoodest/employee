FROM golang:tip-trixie

WORKDIR /employee

COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY go.mod go.sum ./

ENV POSTGRES_USER=postgres
ENV POSTGRES_DB_NAME=postgres
ENV HOST=postgres

RUN go mod download

ENTRYPOINT ["go", "run", "./cmd/main.go"]