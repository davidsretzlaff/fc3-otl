# 🔍 Guia de Queries Loki - Investigação de Logs

## 📚 Queries Básicas

### 1. Ver todos os logs recentes
```logql
{job="containers"}
```

### 2. Filtrar por serviço específico
```logql
{job="containers"} |= "subscription"  # Logs do Subscription Service
{job="containers"} |= "customer"      # Logs do Customer Service
```

### 3. Ver apenas erros
```logql
{job="containers"} |= "ERROR"
{job="containers"} |= "level\":\"error\""  # Para logs JSON estruturados
```

## 🎯 Queries para Investigação

### 4. Buscar por Correlation ID específico
```logql
{job="containers"} |= "subscription-20250604023320-fdf2c569"
{job="containers"} |= "CorrelationId:subscription-20250604023320"
```

### 5. Buscar por cliente específico
```logql
{job="containers"} |= "joaoexemplo@gmail.com"
{job="containers"} |= "customer_email"
```

### 6. Buscar por operação específica
```logql
{job="containers"} |= "ActivateSubscription"
{job="containers"} |= "CreateCustomer"
{job="containers"} |= "operation\":\"activate_subscription\""
```

## ⚡ Queries Avançadas

### 7. Performance - Operações lentas (>3 segundos)
```logql
{job="containers"} |= "duration_ms" | json | duration_ms > 3000
```

### 8. Códigos de status HTTP de erro
```logql
{job="containers"} |= "status code 500"
{job="containers"} |= "status code 4"  # Todos os 4xx
{job="containers"} |= "status code 5"  # Todos os 5xx
```

### 9. Combinações complexas - Erros de um cliente específico
```logql
{job="containers"} |= "ERROR" |= "joaoexemplo@gmail.com"
```

### 10. Rastrear fluxo completo por Correlation ID
```logql
{job="containers"} |= "subscription-20250604023320-fdf2c569" 
| line_format "{{.timestamp}} [{{.service}}] {{.message}}"
```

## 📊 Queries para Métricas

### 11. Contar erros por minuto
```logql
sum(count_over_time({job="containers"} |= "ERROR" [1m])) by (service)
```

### 12. Taxa de erro por serviço
```logql
sum(rate({job="containers"} |= "ERROR" [5m])) by (service) /
sum(rate({job="containers"}[5m])) by (service)
```

### 13. Duração média das operações
```logql
avg_over_time({job="containers"} |= "duration_ms" | json | unwrap duration_ms [5m])
```

## 🚨 Queries para Alertas

### 14. Erros críticos nos últimos 5 minutos
```logql
{job="containers"} |= "ERROR" |= "status code 5" | count_over_time({} [5m]) > 5
```

### 15. Falhas de comunicação entre serviços
```logql
{job="containers"} |= "customer_service_error" |= "connection"
```

## 🔧 Queries para Debug

### 16. Ver logs formatados JSON
```logql
{job="containers"} | json | line_format "{{.timestamp}} [{{.level}}] {{.message}}"
```

### 17. Extrair campos específicos
```logql
{job="containers"} | json | line_format "Correlation: {{.correlation_id}} | Customer: {{.customer_email}} | Error: {{.error}}"
```

### 18. Filtrar por trace ID (OpenTelemetry)
```logql
{job="containers"} |= "trace_id" | json | trace_id =~ ".*abc123.*"
```

## 🎨 Dicas de Uso

### Operadores úteis:
- `|=` : contém
- `!=` : não contém  
- `=~` : regex match
- `!~` : regex não match
- `>`, `<`, `>=`, `<=` : comparações numéricas

### Filtros de tempo:
- `[5m]` : últimos 5 minutos
- `[1h]` : última hora
- `[1d]` : último dia

### Formatação:
- `| json` : parseia JSON
- `| line_format` : formata saída
- `| unwrap` : extrai valores numéricos

## 📝 Exemplos Práticos

### Investigar erro específico:
1. Comece com o correlation ID: `{job="containers"} |= "subscription-20250604023320-fdf2c569"`
2. Veja o contexto: adicione `| line_format "{{.timestamp}} {{.message}}"`
3. Busque logs relacionados: use o customer_email encontrado

### Monitorar performance:
1. Operações lentas: `{job="containers"} |= "duration_ms" | json | duration_ms > 3000`
2. Por serviço: adicione `|= "subscription"` ou `|= "customer"`
3. Timeline: use `line_format` para ver a sequência temporal

### Debug de fluxo completo:
1. Pegue o correlation ID do erro
2. Use: `{job="containers"} |= "SEU_CORRELATION_ID" | line_format "{{.timestamp}} [{{.service}}] {{.operation}} {{.message}}"`
3. Ordene por timestamp para ver o fluxo 