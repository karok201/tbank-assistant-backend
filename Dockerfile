FROM golang:onbuild

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
#USER www-data

COPY . /app

FROM alpine:latest

COPY main.go go.mod go.sum ./

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o server

ENTRYPOINT ["/entrypoint.sh"]
