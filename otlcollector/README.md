# OTL Collector

Este é o serviço de coleta de telemetria que atua como intermediário entre as aplicações e o Jaeger.

## Configuração

O arquivo `config.yaml` define três componentes principais:

1. **Receivers**: Configurado para receber dados OTLP via HTTP na porta 4318
2. **Processors**: Remove campos sensíveis dos spans (como cartões de crédito, senhas e tokens)
3. **Exporters**: Envia os dados processados para o Jaeger via gRPC

## Pipeline de Dados

1. As aplicações (Go e .NET) enviam dados para o collector via HTTP
2. O collector processa os dados, removendo informações sensíveis
3. Os dados processados são enviados para o Jaeger via gRPC

## Como Usar

O collector é iniciado automaticamente com o docker-compose. As aplicações devem configurar o endpoint do collector como:

```
OTEL_EXPORTER_OTLP_ENDPOINT=http://otlcollector:4318
```

## Monitoramento

Os dados podem ser visualizados no Jaeger UI em:
http://localhost:16686 