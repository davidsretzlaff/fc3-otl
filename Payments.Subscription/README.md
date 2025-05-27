# Payments Customer API

API para gerenciamento de customers com arquitetura baseada em DDD (Domain Driven Design) e OpenTelemetry para observabilidade.

## 🚀 Tecnologias

- **Go 1.21+**
- **MySQL 8.0**
- **OpenTelemetry** para tracing
- **Gorilla Mux** para roteamento HTTP
- **Docker Compose** para ambiente de desenvolvimento

## 📋 Pré-requisitos

- Go 1.21 ou superior
- Docker e Docker Compose
- Git

## 🛠️ Setup do Projeto

### 1. Clone o repositório
```bash
git clone <seu-repositorio>
cd Payments.Customer
```

### 2. Instale as dependências
```bash
go mod tidy
```

### 3. Inicie o banco de dados MySQL
```bash
docker-compose up -d
```

### 4. Execute a aplicação
```bash
go run cmd/api/main.go
```

A API estará disponível em: `http://localhost:8888`

## 📚 Endpoints da API

### Health Check

#### GET /health
Verifica se a API está funcionando.

**cURL:**
```bash
curl -X GET http://localhost:8888/health
```

**Postman:**
- Method: `GET`
- URL: `http://localhost:8888/health`

**Resposta:**
```json
{
  "status": "ok"
}
```

---

### Criar Customer

#### POST /customers
Cria um novo customer no sistema.

**cURL:**
```bash
curl -X POST http://localhost:8888/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao.silva@email.com",
    "document": "12345678901"
  }'
```

**Postman:**
- Method: `POST`
- URL: `http://localhost:8888/customers`
- Headers:
  - `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "name": "João Silva",
  "email": "joao.silva@email.com",
  "document": "12345678901"
}
```

**Resposta de Sucesso (201):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "João Silva",
  "email": "joao.silva@email.com",
  "document": "12345678901",
  "status": "active",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

---

### Buscar Customer por ID

#### GET /customers/{id}
Busca um customer específico pelo ID.

**cURL:**
```bash
curl -X GET http://localhost:8888/customers/550e8400-e29b-41d4-a716-446655440000
```

**Postman:**
- Method: `GET`
- URL: `http://localhost:8888/customers/550e8400-e29b-41d4-a716-446655440000`

**Resposta de Sucesso (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "João Silva",
  "email": "joao.silva@email.com",
  "document": "12345678901",
  "status": "active",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Resposta de Erro (404):**
```json
{
  "error": "customer não encontrado"
}
```

---

### Listar Todos os Customers

#### GET /customers
Retorna uma lista com todos os customers cadastrados, ordenados por data de criação (mais recentes primeiro).

**cURL:**
```bash
curl -X GET http://localhost:8888/customers
```

**Postman:**
- Method: `GET`
- URL: `http://localhost:8888/customers`

**Resposta de Sucesso (200):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "João Silva",
    "email": "joao.silva@email.com",
    "document": "12345678901",
    "status": "active",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": "660f9511-f3ac-52e5-b827-557766551111",
    "name": "Maria Santos",
    "email": "maria.santos@email.com",
    "document": "98765432100",
    "status": "active",
    "created_at": "2024-01-15T09:15:00Z",
    "updated_at": "2024-01-15T09:15:00Z"
  }
]
```

**Resposta quando não há customers (200):**
```json
[]
```

---

## 🧪 Testando a API

### Sequência de Testes Completa

1. **Verificar se a API está funcionando:**
```bash
curl -X GET http://localhost:8888/health
```

2. **Criar alguns customers:**
```bash
# Customer 1
curl -X POST http://localhost:8888/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao.silva@email.com",
    "document": "12345678901"
  }'

# Customer 2
curl -X POST http://localhost:8888/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Maria Santos",
    "email": "maria.santos@email.com",
    "document": "98765432100"
  }'
```

3. **Listar todos os customers para ver os IDs:**
```bash
curl -X GET http://localhost:8888/customers
```

4. **Buscar um customer específico (use um ID da lista anterior):**
```bash
curl -X GET http://localhost:8888/customers/{ID_DO_CUSTOMER}
```

---

## 🏗️ Arquitetura

O projeto segue os princípios de **Domain Driven Design (DDD)** e **SOLID**:

```
internal/customer/
├── customer.go          # Domain (Aggregate + Repository Interface)
├── service.go           # Application Service
├── handler.go           # HTTP Handlers
└── mysql/repository.go  # Infrastructure (MySQL Repository)
```

### Camadas:

- **Domain**: Entidades, Value Objects e regras de negócio
- **Application**: Casos de uso e orquestração
- **Infrastructure**: Implementações de repositórios e integrações
- **Presentation**: Controllers HTTP e serialização

---

## 🐳 Docker

### Banco de Dados MySQL

O projeto inclui um `docker-compose.yml` que configura:

- **MySQL 8.0** na porta `3306`
- **Database**: `payments_db`
- **User**: `payments` / **Password**: `payments123`
- **Tabela `customers`** criada automaticamente

### Comandos Docker úteis:

```bash
# Iniciar o banco
docker-compose up -d

# Ver logs do banco
docker-compose logs mysql

# Parar o banco
docker-compose down

# Resetar o banco (apaga todos os dados)
docker-compose down -v && docker-compose up -d
```

---

## 📊 Observabilidade

A aplicação inclui **OpenTelemetry** para tracing distribuído. Os traces são gerados automaticamente para:

- Requisições HTTP
- Operações de banco de dados
- Operações de serviço

---

## 🔧 Configuração

### Configurações da Aplicação:
- **Porta**: `8888`
- **Host do MySQL**: `localhost:3306`
- **Database**: `payments_db`
- **User**: `payments`
- **Password**: `payments123`

### Variáveis de Ambiente (Opcionais):
Atualmente as configurações estão fixas no código, mas podem ser facilmente migradas para variáveis de ambiente se necessário.

---

## 🚨 Códigos de Status HTTP

| Status | Descrição |
|--------|-----------|
| 200 | Sucesso |
| 201 | Criado com sucesso |
| 400 | Dados inválidos |
| 404 | Recurso não encontrado |
| 500 | Erro interno do servidor |

---

## 📝 Exemplos de Erro

### Dados inválidos (400):
```json
{
  "error": "nome do customer é obrigatório"
}
```

### Customer não encontrado (404):
```json
{
  "error": "customer não encontrado"
}
```

### Erro interno (500):
```json
{
  "error": "erro ao conectar com o banco de dados"
}
```

---

## 🤝 Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

---

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes. 