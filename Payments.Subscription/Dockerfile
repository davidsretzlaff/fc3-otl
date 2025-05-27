# Estágio de build
FROM golang:1.23-alpine AS builder

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos de dependência primeiro
COPY ./go.mod ./go.sum ./

# Baixa as dependências
RUN go mod download

# Copia todo o código-fonte do projeto
COPY . .

# Compila o código-fonte principal localizado em cmd/api/
RUN CGO_ENABLED=0 GOOS=linux go build -o payments.subscription ./cmd/api

# Estágio final
FROM alpine:latest

# Instala certificados (para HTTPS funcionar)
RUN apk --no-cache add ca-certificates

# Define o diretório de trabalho
WORKDIR /app

# Copia o binário compilado
COPY --from=builder /app/payments.subscription .

# Expõe a porta da aplicação
EXPOSE 8888

# Comando de inicialização
CMD ["./payments.subscription"]
