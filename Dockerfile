FROM golang:1.18.0-buster AS build

ENV SERVE_ON_PORT="8000"

WORKDIR /app

RUN ls
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /daemon-exercise ./cmd/main.go

