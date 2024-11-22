FROM golang:onbuild

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
#USER www-data

COPY . /app

FROM alpine:latest

RUN apk add --no-cache git

RUN go get github.com/gin-gonic/gin
RUN go get github.com/golang-jwt/jwt/v5
RUN go get github.com/joho/godotenv
RUN go get golang.org/x/crypto
RUN go get gorm.io/driver/mysql
RUN go get gorm.io/gorm
RUN go get github.com/githubnemo/CompileDaemon
RUN go get github.com/spf13/viper
RUN go get gorm.io/gorm

ENTRYPOINT ["/entrypoint.sh"]
