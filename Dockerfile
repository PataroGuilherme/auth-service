# ---------------------------
# Etapa 1: Build da Aplicação
# ---------------------------
FROM golang:1.22 AS builder

# Diretório de trabalho
WORKDIR /app

# Copia arquivos de dependência
COPY go.mod ./
RUN go mod download

# Copia todo o código
COPY . .

# Compila o binário estático
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service main.go


# --------------------------------
# Etapa 2: Imagem Final (Alpine)
# --------------------------------
FROM alpine:3.20

# Diretório da aplicação
WORKDIR /app

# Copia binário da etapa de build
COPY --from=builder /app/auth-service .

# Expõe as portas solicitadas
EXPOSE 8000
EXPOSE 8080
EXPOSE 80
EXPOSE 443

# Comando padrão
CMD ["./auth-service"]
