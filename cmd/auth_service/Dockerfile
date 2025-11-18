FROM golang:1.25.1-alpine

RUN apk add git bash curl && \
    go install github.com/swaggo/swag/cmd/swag@latest 

# Настраиваем окружение Go
RUN go env -w GOPROXY=https://proxy.golang.org,direct && \
    go env -w GOSUMDB=off

WORKDIR /app

COPY ./go.* .
RUN go mod download

COPY . .
RUN go build -o /app/main ./cmd/auth_service/main.go 

EXPOSE 8084
CMD ["/app/main"]