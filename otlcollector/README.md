# OpenTelemetry Collector - Configuração para Logs

## Visão Geral

Este OpenTelemetry Collector foi configurado para **substituir completamente o Loki + Promtail**, processando logs diretamente e enviando para o Loki com labels apropriados para filtros no Grafana.

## Funcionalidades

### 📋 **Processamento de Logs JSON**
- Processa logs no formato: `{"time":"2025-06-07T15:39:57Z","level":"information","msg":"[customer] Saving 1 customers to database","correlation_id":"subscription-20250607153957-5d9deb23","service":"customer"}`
- Extrai automaticamente campos: `time`, `level`, `msg`, `correlation_id`, `service`
- Converte campos em labels para filtragem no Grafana

### 🏷️ **Labels Criados**
Os mesmos labels que o Loki/Promtail criava:
- **service**: Nome do serviço (customer, payment, subscription)
- **level**: Nível do log (info, error, warn, debug)
- **correlation_id**: ID de correlação para rastreamento
- **environment**: Ambiente de deployment (development)

### 🔄 **Normalização de Níveis**
- `information` → `info`
- `warning` → `warn`
- `critical` → `error`
- `debug` → `debug`

## Configuração

### Receivers
- **OTLP**: Recebe traces/metrics/logs via HTTP (4318) e gRPC (4317)
- **FileLog**: Monitora arquivos `/logs/*.log` para logs JSON

### Processors
- **transform/add_labels**: Adiciona labels baseados nos campos JSON
- **resource**: Adiciona resource attributes (environment)
- **batch**: Agrupa logs para envio eficiente

### Exporters
- **loki**: Envia logs diretamente para Loki com labels
- **otlp**: Envia traces para Jaeger
- **prometheus**: Expõe métricas
- **logging**: Debug logs no console

## Migração do Loki/Promtail

### ✅ **Removido**
- Promtail container
- Dependência Promtail → Loki
- Configuração manual de labels

### ✅ **Adicionado**
- Processamento nativo de logs JSON
- Labels automáticos baseados em campos
- Envio direto para Loki
- Normalização de níveis de log

## Uso no Grafana

Os mesmos filtros que funcionavam com Loki/Promtail continuam funcionando:

```logql
# Filtrar por serviço
{service="customer"}

# Filtrar por nível
{level="error"}

# Filtrar por correlation_id
{correlation_id="subscription-20250607153957-5d9deb23"}

# Combinações
{service="customer", level="error"}
```

## Estrutura do Log

Mantenha o formato JSON atual:
```json
{
  "time": "2025-06-07T15:39:57Z",
  "level": "information",
  "msg": "[customer] Saving 1 customers to database",
  "correlation_id": "subscription-20250607153957-5d9deb23",
  "service": "customer"
}
```

## Deployment

```bash
# Rebuild do OpenTelemetry Collector
docker-compose build otlcollector

# Restart dos serviços (sem Promtail)
docker-compose up -d
``` 