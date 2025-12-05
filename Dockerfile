FROM golang:1.22 AS builder

WORKDIR /app

# Copia SOMENTE os arquivos de dependências primeiro
COPY go.mod go.sum ./
RUN go mod download

# Agora copia TODOS os arquivos Go do serviço
COPY *.go ./

# Aqui rodamos o build
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/auth-service .

EXPOSE 8000 8080 80 443

CMD ["./auth-service"]
