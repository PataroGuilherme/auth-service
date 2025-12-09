# ============================
# 1 - Build da aplicação
# ============================
FROM golang:1.22 AS builder

WORKDIR /app

# Copia go.mod e go.sum primeiro para aproveitar cache
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service .

# ============================
# 2 - Imagem final
# ============================
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/auth-service /app/auth-service

# Porta do serviço
EXPOSE 8001

# Inicia o serviço
CMD ["/app/auth-service"]
