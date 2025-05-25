# Projeto FC3-OTL

Este projeto segue a arquitetura de camadas baseada em Domain-Driven Design (DDD) e princípios SOLID.

## 📁 Estrutura do Projeto

```
netcoreapp/
├── app-otl/                    # 🌐 API - Camada de apresentação
│   ├── Api.csproj             # Projeto Web API
│   └── ...
├── src/
│   ├── Domain/                # 🏗️ DOMAIN - Camada de domínio
│   │   ├── Domain.csproj      # Entidades, Value Objects, Aggregates
│   │   └── Domain/
│   ├── Application/           # 📋 APPLICATION - Camada de aplicação
│   │   ├── Application.csproj # Use Cases, Services, DTOs
│   │   └── UseCases/
│   └── Infra.Data/           # 💾 INFRASTRUCTURE - Camada de infraestrutura
│       ├── Infra.Data.csproj # Repositórios, DbContext, Dados
│       └── AppDbContext.cs
└── netcoreapp.sln            # Solution file
```

## 🏗️ Arquitetura de Camadas

### 🏗️ Domain (Domínio)
- **Responsabilidade**: Lógica de negócio central
- **Conteúdo**: Entidades, Value Objects, Aggregates, Domain Services
- **Dependências**: Nenhuma (camada mais interna)

### 📋 Application (Aplicação)
- **Responsabilidade**: Casos de uso e orquestração
- **Conteúdo**: Use Cases, Application Services, DTOs, Interfaces
- **Dependências**: Domain

### 💾 Infra.Data (Infraestrutura de Dados)
- **Responsabilidade**: Persistência e acesso a dados
- **Conteúdo**: Repositórios, DbContext, Migrations, Configurações EF
- **Dependências**: Domain, Application
- **Tecnologias**: Entity Framework Core, SQL Server

### 🌐 API (Apresentação)
- **Responsabilidade**: Interface externa e controladores
- **Conteúdo**: Controllers, Middlewares, Configurações
- **Dependências**: Domain, Application, Infra.Data
- **Tecnologias**: ASP.NET Core, OpenTelemetry, Serilog

## 🛠️ Comandos Úteis

### Build da Solution
```bash
dotnet build netcoreapp.sln
```

### Restaurar Pacotes
```bash
dotnet restore netcoreapp.sln
```

### Executar API
```bash
cd app-otl
dotnet run
```

### Executar Testes
```bash
dotnet test
```

## 🔗 Dependências Entre Projetos

```
API → Application → Domain
API → Infra.Data → Domain
API → Infra.Data → Application
```

## 📦 Principais Pacotes

- **Entity Framework Core** - ORM para acesso a dados
- **Dapper** - Micro ORM para queries otimizadas
- **OpenTelemetry** - Observabilidade e telemetria
- **Serilog** - Logging estruturado

## 🚀 Getting Started

1. **Clone o repositório**
2. **Restaure os pacotes**: `dotnet restore`
3. **Build a solution**: `dotnet build`
4. **Execute a API**: `cd app-otl && dotnet run`

## 📋 Próximos Passos

- [ ] Configurar Entity Framework Migrations
- [ ] Implementar padrão Repository
- [ ] Adicionar testes unitários
- [ ] Configurar Docker Compose
- [ ] Implementar autenticação JWT 