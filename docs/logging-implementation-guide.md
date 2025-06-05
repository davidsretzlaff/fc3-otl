# 🎯 Guia de Implementação de Logs Profissionais com Correlation ID

## ✅ **O que foi implementado:**

### **1. Sistema de Correlation ID**
- **Geração automática** no primeiro serviço (Subscription Service)
- **Propagação via headers HTTP** (`X-Correlation-ID`)
- **Integração com OpenTelemetry** tracing
- **Formato padronizado**: `{service}-{timestamp}-{random}`

### **2. Logging Estruturado (Go - Subscription Service)**
```json
{
  "timestamp": "2024-01-15T10:30:00.123Z",
  "level": "INFO",
  "service": "subscription-service",
  "correlation_id": "subscription-20241115-abc123def456",
  "trace_id": "47a4d9efca74d6b962c64bd5a0d4f83d",
  "span_id": "c42a0b27839aa914",
  "operation": "CreateSubscription",
  "duration_ms": 342,
  "message": "Customer criado com sucesso",
  "context": {
    "customer_id": "cust-123",
    "customer_email": "user@example.com"
  }
}
```

### **3. Logging Estruturado (C# - Customer Service)**
```
[2024-01-15 10:30:00.123 +00:00] [INF] [Customer.API] [CorrelationId:subscription-20241115-abc123def456] [TraceId:47a4d9efca74d6b962c64bd5a0d4f83d] [SpanId:c42a0b27839aa914] Customer created successfully. CustomerId: cust-123, Duration: 234ms
```

## 🧪 **Como Testar:**

### **1. Teste Básico de Correlation ID**
```bash
# 1. Criar uma subscription
curl -X POST http://localhost:8081/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "plan_id": "basic-plan",
    "customer": {
      "name": "João Silva",
      "email": "joao@example.com"
    }
  }'

# 2. Observar nos logs que o mesmo correlation_id aparece em ambos os serviços
```

### **2. Teste com Correlation ID Personalizado**
```bash
# Enviar requisição com correlation ID específico
curl -X POST http://localhost:8081/subscriptions \
  -H "Content-Type: application/json" \
  -H "X-Correlation-ID: test-custom-correlation-123" \
  -d '{
    "plan_id": "premium-plan",
    "customer": {
      "name": "Maria Santos",
      "email": "maria@example.com"
    }
  }'
```

### **3. Verificação dos Logs**
```bash
# Subscription Service (Go)
docker logs fc3-otl-payments.subscription-1 | grep "test-custom-correlation-123"

# Customer Service (C#)
docker logs fc3-otl-payments.customer-1 | grep "test-custom-correlation-123"
```

## 📊 **Campos de Log Implementados:**

### **Campos Obrigatórios (Todos os logs):**
- ✅ `timestamp` - Timestamp UTC com precisão de milissegundos
- ✅ `level` - Nível do log (INFO, WARN, ERROR, etc.)
- ✅ `service` - Nome do serviço
- ✅ `correlation_id` - ID de correlação único por transação
- ✅ `operation` - Nome da operação sendo executada
- ✅ `message` - Mensagem descritiva

### **Campos Condicionais:**
- ✅ `trace_id` - ID do trace do OpenTelemetry
- ✅ `span_id` - ID do span atual
- ✅ `duration_ms` - Duração da operação (quando aplicável)
- ✅ `error` - Detalhes do erro (apenas em logs de erro)
- ✅ `context` - Dados específicos da operação

### **Logs de Operações Externas:**
- ✅ `target_service` - Serviço de destino
- ✅ `status_code` - Código de status HTTP
- ✅ `url` - URL da requisição

## 🎯 **Principais Funcionalidades:**

### **1. Alertas Automáticos**
- ⚠️ **Operações lentas** (> 3 segundos): Log de WARNING
- 🚨 **Operações muito lentas** (> 5 segundos): Log de ERROR
- 📊 **Métricas de duração** em todos os logs de fim de operação

### **2. Contexto Completo**
- 🔄 **Trace completo** da requisição através dos serviços
- 🏷️ **Tags de negócio** (customer_email, plan_id, etc.)
- 🔍 **Facilita debugging** com correlation ID

### **3. Estrutura Profissional**
- 📝 **JSON estruturado** (Go) e template personalizado (C#)
- 🔒 **Sem dados sensíveis** nos logs
- 📈 **Compatível com Grafana/ELK Stack**

## 🔍 **Pontos de Observação:**

### **1. Fluxo Completo da Requisição**
```
[Subscription Service] correlation_id: sub-20241115-abc123
├── Operation: CreateSubscription started
├── Operation: CustomerClient.CreateCustomer started  
├── External call to customer-service                    ← Headers propagados
│   
[Customer Service] correlation_id: sub-20241115-abc123   ← Mesmo correlation_id!
├── Operation: CreateCustomer started
├── Customer created successfully
│
[Subscription Service] correlation_id: sub-20241115-abc123
├── Customer criado com sucesso
├── Subscription created successfully
└── Operation: CreateSubscription completed (duration: 342ms)
```

### **2. Logs de Erro com Contexto**
```json
{
  "level": "ERROR",
  "correlation_id": "sub-20241115-abc123",
  "operation": "CreateCustomer", 
  "message": "Erro ao criar customer",
  "error": "connection timeout",
  "context": {
    "customer_email": "joao@example.com",
    "plan_id": "basic-plan",
    "target_service": "customer-service"
  }
}
```

## 🚀 **Próximos Passos (Opcionais):**

### **1. Integração com Grafana**
- Dashboard com correlation_id como filtro
- Alertas baseados em logs de ERROR
- Métricas de performance por operação

### **2. Centralização de Logs**
- Envio para ELK Stack ou Loki
- Índices por correlation_id para busca rápida
- Retenção de logs baseada em criticidade

### **3. Monitoramento Proativo**
- Alertas Slack/Teams para erros críticos
- SLA tracking baseado em duration_ms
- Health checks com correlation ID

## 📋 **Checklist de Validação:**

- [ ] Correlation ID é gerado automaticamente
- [ ] Correlation ID é propagado entre serviços
- [ ] Logs incluem trace_id e span_id
- [ ] Operações lentas geram warnings
- [ ] Erros incluem contexto completo
- [ ] Headers HTTP são preservados
- [ ] Formato JSON é válido (Go)
- [ ] Template Serilog funciona (C#)
- [ ] Health checks funcionam
- [ ] Logs são enviados para OpenTelemetry Collector

✅ **Sistema de logs profissional implementado com sucesso!**