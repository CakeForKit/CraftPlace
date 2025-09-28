FROM golang:1.25.1-alpine

RUN apk add git bash && \
    go install github.com/swaggo/swag/cmd/swag@latest 
RUN go install github.com/air-verse/air@latest

# Настраиваем окружение Go
RUN go env -w GOPROXY=https://proxy.golang.org,direct && \
    go env -w GOSUMDB=off

WORKDIR /app

# Копируем только файлы модулей сначала для эффективного кэширования
COPY ./go.* .
RUN go mod download

RUN mkdir -p tmp
CMD ["air", "-c", "./cmd/dev/.air.toml"]