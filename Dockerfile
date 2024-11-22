FROM golang:alpine AS builder

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY . /app

FROM alpine:latest

COPY main.go go.mod go.sum ./

ADD go.mod .

ADD go.sum .

RUN go mod download

ENTRYPOINT ["/entrypoint.sh"]
