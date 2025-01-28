FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN ./generate-api-docs.sh

RUN go build -o url-shortener

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/url-shortener /

EXPOSE 8080

CMD ["/url-shortener"]

