FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/auth-service .

# Variáveis que serão fornecidas pelo docker-compose
ENV DB_HOST=auth-db
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASS=postgres
ENV DB_NAME=authdb

EXPOSE 8000 8080 80 443

CMD ["./auth-service"]
