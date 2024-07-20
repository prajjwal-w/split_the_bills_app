# Use the official Golang image as the base image
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY .env .env

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]