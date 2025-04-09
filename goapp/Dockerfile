# Estágio de build
FROM golang:1.21-alpine AS builder

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos go.mod e go.sum primeiro (melhor prática para cache de camadas)
COPY go.mod go.sum ./

# Baixa as dependências do projeto
RUN go mod download

# Copia todo o código fonte
COPY . .

# Compila a aplicação
# CGO_ENABLED=0: Desativa o CGO para criar um binário estático
# -o goapp: Nome do binário de saída
RUN CGO_ENABLED=0 GOOS=linux go build -o goapp

# Estágio final (imagem menor)
FROM alpine:latest

# Instala o ca-certificates (necessário para HTTPS)
RUN apk --no-cache add ca-certificates

# Define o diretório de trabalho
WORKDIR /app

# Copia o binário compilado do estágio anterior
COPY --from=builder /app/goapp .

# Expõe a porta que a aplicação usa
EXPOSE 8888

# Comando que será executado quando o container iniciar
CMD ["./goapp"]

