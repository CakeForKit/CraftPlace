FROM golang:1.25.1-alpine

RUN apk add git bash
RUN go install github.com/air-verse/air@latest

# Настраиваем окружение Go
RUN go env -w GOPROXY=https://proxy.golang.org,direct && \
    go env -w GOSUMDB=off

WORKDIR /app

COPY ./go.* .
RUN go mod download
COPY . .
RUN go build -o /app/main ./cmd/replica/main.go 


RUN ls -la /app
EXPOSE 8081
CMD ["/app/main"]