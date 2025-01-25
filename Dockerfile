FROM golang:1.23-bookworm AS base

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o url-shortener

EXPOSE 8080

RUN ./generate-api-docs.sh

CMD ["./url-shortener"]

