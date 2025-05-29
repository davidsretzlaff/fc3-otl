# Tracing Distribuído e Propagação de Contexto

## Conceitos Básicos

### Trace
Um `trace` representa a jornada completa de uma requisição através de um sistema distribuído. Pense nele como o "caminho" que uma requisição percorre através de diferentes serviços.

### Span
Um `span` representa uma única operação ou unidade de trabalho dentro de um trace. Por exemplo:
- Uma chamada HTTP
- Uma query no banco de dados
- Um processamento específico

### IDs Importantes

#### TraceID
- É um identificador único para todo o trace
- Permanece o mesmo em todos os serviços que a requisição passa
- Permite visualizar a jornada completa da requisição
- Exemplo: `47a4d9efca74d6b962c64bd5a0d4f83d`

#### SpanID
- Identifica unicamente um span específico
- Cada operação dentro do trace tem seu próprio SpanID
- Exemplo: `2d281dd73dcf7574`

#### ParentSpanID
- Referencia o SpanID do span que originou a operação atual
- Permite construir a hierarquia/árvore de spans
- Exemplo: `c42a0b27839aa914`

## Propagação de Contexto

### Como Funciona
1. **Serviço A (Origem)**
   ```go
   // Go (usando OpenTelemetry)
   ctx, span := tracer.Start(ctx, "OperacaoA")
   defer span.End()
   
   // Injeta o contexto nos headers HTTP
   propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
   ```

2. **Headers HTTP**
   ```
   traceparent: 00-47a4d9efca74d6b962c64bd5a0d4f83d-c42a0b27839aa914-01
   ```
   Formato: `versão-traceID-spanID-flags`

3. **Serviço B (Destino)**
   ```csharp
   // C# (usando OpenTelemetry)
   // O ASP.NET Core extrai automaticamente o contexto dos headers
   app.UseOpenTelemetry();
   ```

### Headers de Propagação

#### traceparent
- Carrega as informações essenciais do trace
- Formato: `versão-traceID-spanID-flags`
- Exemplo: `00-47a4d9efca74d6b962c64bd5a0d4f83d-c42a0b27839aa914-01`

#### tracestate
- Carrega informações adicionais do trace
- Usado para dados específicos do vendor
- Opcional

## CorrelationID vs TraceID

### CorrelationID
- Identificador único para correlacionar logs e eventos
- Geralmente mais simples, apenas um ID
- Frequentemente usado em logs estruturados
- Não carrega informações de hierarquia

### TraceID
- Parte de um sistema mais complexo de tracing
- Carrega informações de hierarquia (ParentSpanID)
- Permite visualização detalhada do fluxo
- Integrado com ferramentas de APM (Application Performance Monitoring)

### Quando Usar Cada Um?
- **CorrelationID**: Para sistemas mais simples, onde só precisa correlacionar logs
- **TraceID**: Para sistemas distribuídos complexos, onde precisa entender o fluxo completo

## Exemplo Prático

### Fluxo Subscription -> Customer

1. **Subscription Inicia o Trace**
   ```go
   ctx, span := tracer.Start(ctx, "CreateSubscription")
   ```

2. **Subscription Chama Customer**
   ```go
   // Cria um novo span filho
   ctx, span := tracer.Start(ctx, "CustomerClient.CreateCustomer")
   
   // Headers são injetados automaticamente
   propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
   ```

3. **Customer Recebe a Requisição**
   ```csharp
   // ASP.NET Core extrai o contexto automaticamente
   // O novo span terá:
   // - Mesmo TraceID do Subscription
   // - Novo SpanID próprio
   // - ParentSpanID referenciando o span do Subscription
   ```

## Visualização no Jaeger

```
Trace
├── Subscription: /subscriptions (SpanID: abc123)
│   └── CustomerClient.CreateCustomer (SpanID: def456)
│       └── Customer: POST /api/customer (SpanID: ghi789)
│           └── Customer: Processamento (SpanID: jkl012)
```

## Benefícios

1. **Debugging Distribuído**
   - Visualize o fluxo completo da requisição
   - Identifique gargalos e falhas
   - Entenda dependências entre serviços

2. **Performance**
   - Meça tempos de resposta
   - Identifique operações lentas
   - Otimize baseado em dados reais

3. **Observabilidade**
   - Monitore sistemas distribuídos
   - Entenda comportamentos em produção
   - Facilite a resolução de problemas

## Ferramentas Populares

1. **OpenTelemetry**
   - Padrão open-source para instrumentação
   - Suporte multi-linguagem
   - Exportadores para diversos backends

2. **Jaeger**
   - Interface visual para traces
   - Análise de performance
   - Filtros e buscas avançadas

3. **Zipkin**
   - Alternativa ao Jaeger
   - Foco em baixa latência
   - Boa integração com Spring 

## Debugging e Troubleshooting

### Verificando Headers e Contexto

1. **Logs no Serviço de Origem (Go)**
   ```go
   // Adicione logs antes de enviar a requisição
   fmt.Printf("Headers sendo enviados: %+v\n", req.Header)
   ```

2. **Logs no Serviço de Destino (C#)**
   ```csharp
   app.Use(async (context, next) =>
   {
       var activity = System.Diagnostics.Activity.Current;
       if (activity != null)
       {
           // Log dos headers de tracing
           var traceparent = context.Request.Headers["traceparent"].ToString();
           var tracestate = context.Request.Headers["tracestate"].ToString();
           
           Log.Information("Trace Headers - TraceParent: {TraceParent}, TraceState: {TraceState}", 
               traceparent, 
               tracestate);
           
           Log.Information("Activity Info - TraceId: {TraceId}, SpanId: {SpanId}, ParentSpanId: {ParentSpanId}",
               activity.TraceId,
               activity.SpanId,
               activity.ParentSpanId);
       }
       await next();
   });
   ```

### Comandos Docker para Debug

1. **Verificar logs de um serviço específico**
   ```bash
   # Ver logs do serviço customer
   docker-compose logs payments.customer

   # Ver logs do serviço subscription
   docker-compose logs payments.subscription

   # Ver logs em tempo real
   docker-compose logs -f payments.customer
   ```

2. **Verificar logs dentro do container**
   ```bash
   # Acessar logs dentro do container
   docker exec fc3-otl-payments.customer-1 cat /app/logs/app.log
   ```

3. **Reconstruir containers após mudanças**
   ```bash
   # Reconstruir um serviço específico
   docker-compose up -d --build payments.customer
   docker-compose up -d --build payments.subscription

   # Reconstruir todos os serviços
   docker-compose up -d --build
   ```

### O que Procurar nos Logs

1. **No Serviço de Origem**
   - Headers `traceparent` sendo injetados
   - TraceID gerado
   - SpanID do span atual

2. **No Serviço de Destino**
   - Headers `traceparent` recebidos
   - Mesmo TraceID do serviço de origem
   - ParentSpanID igual ao SpanID do span de origem
   - Novo SpanID gerado

### Problemas Comuns

1. **Headers Não Chegando**
   - Verifique se o propagator está configurado corretamente
   - Confirme se a injeção dos headers está acontecendo
   - Verifique se não há proxy/gateway removendo os headers

2. **TraceID Diferente**
   - Verifique se o contexto está sendo propagado corretamente
   - Confirme se não há múltiplas configurações do OpenTelemetry

3. **Spans Não Aparecendo no Jaeger**
   - Verifique a configuração do coletor OpenTelemetry
   - Confirme se os spans estão sendo finalizados (defer span.End())
   - Verifique se o exportador está configurado corretamente

### Dicas de Debugging

1. **Use Tags nos Spans**
   ```go
   span.SetAttributes(
       attribute.String("customer.name", request.Name),
       attribute.String("customer.email", request.Email),
   )
   ```

2. **Adicione Eventos nos Spans**
   ```go
   span.AddEvent("iniciando_chamada_http")
   // ... código ...
   span.AddEvent("chamada_http_completada")
   ```

3. **Capture Erros**
   ```go
   if err != nil {
       span.RecordError(err)
       span.SetStatus(codes.Error, err.Error())
   }
   ``` 