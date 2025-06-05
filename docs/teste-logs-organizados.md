# 🧪 Teste da Nova Configuração - Logs Organizados

## **✅ Passo a Passo para Testar**

### **1. Aplicar a Nova Configuração:**

```bash
# Reiniciar apenas o Promtail com a nova configuração
docker-compose restart promtail

# Aguardar 30 segundos para os logs começarem a aparecer
```

### **2. Testar no Grafana:**

1. **Acesse:** `http://localhost:3000` (admin/admin)
2. **Vá para:** Explore → Loki
3. **Execute essas queries:**

#### **A) Testar Logs de Aplicação (limpos):**
```logql
{job="app"} | limit 10
```

**✅ Resultado esperado (formato real):**
```
13:24:50 [subscription] [CorrelationId:subscription-20250605132450-f8e767a0] Starting CreateSubscription for abaaac@gmail.com
13:24:50 [subscription] [CorrelationId:subscription-20250605132450-f8e767a0] Iniciando criação de customer
13:24:50 [subscription] [CorrelationId:subscription-20250605132450-f8e767a0] Calling customer service to create customer
13:24:50 [customer] [CorrelationId:subscription-20250605132450-f8e767a0] Received create customer request
13:24:50 [customer] [customer] Customer added to in-memory collection. Email: abaaac@gmail.com
13:24:50 [customer] [customer] Saving 1 customers to database
```

#### **B) Testar Logs de Sistema:**
```logql
{job="system"} | limit 10
```

**✅ Resultado esperado:**
```
2025-01-10 [mysql] Ready for connections
2025-01-10 [prometheus] Server is ready
```

#### **C) Testar Logs do Loki:**
```logql
{job="loki"} | limit 10
```

**✅ Resultado esperado:**
```
level=info ts=2025-01-10 msg="Loki started"
```

#### **D) Verificar Separação:**
```logql
{job="app"} |= "mysql"
```

**✅ Resultado esperado:** NENHUM resultado (logs do MySQL não aparecem em app)

### **3. Testes Funcionais:**

#### **A) Gerar Logs de Teste (PowerShell):**

```powershell
# Fazer uma requisição para gerar logs da aplicação
Invoke-RestMethod -Uri "http://localhost:8888/subscription" -Method Post -ContentType "application/json" -Body '{"customer_email": "teste.grafana@teste.com", "plan_id": 1}'
```

#### **B) Verificar no Grafana:**
```logql
{job="app"} |= "teste.grafana@teste.com" | limit 20
```

**✅ Deveria mostrar:**
```
15:30:45 [subscription] [CorrelationId:subscription-...] Starting CreateSubscription for teste.grafana@teste.com
15:30:45 [subscription] [CorrelationId:subscription-...] Calling customer service to create customer
15:30:45 [customer] [CorrelationId:subscription-...] Received create customer request
```

### **4. Testes de Filtros:**

#### **A) Apenas Erros da Aplicação:**
```logql
{job="app"} |= "ERROR" | limit 10
```

#### **B) Por Serviço Específico (usando labels):**
```logql
{job="app", app_service="subscription"} | limit 10
{job="app", app_service="customer"} | limit 10
```

#### **C) Por Correlation ID (usando labels):**
```logql
{job="app", correlation_id=~"subscription-.*"} | limit 30
```

#### **D) Filtros por texto:**
```logql
{job="app"} |= "subscription" | limit 20
{job="app"} |= "customer" | limit 20
{job="app"} |= "CorrelationId" | limit 30
```

## **🎯 O Que Você DEVE Ver:**

### **✅ CERTO - Com {job="app"}:**
```
13:24:50 [subscription] [CorrelationId:subscription-20250605132450-f8e767a0] Starting CreateSubscription for abaaac@gmail.com
13:24:50 [subscription] [CorrelationId:subscription-20250605132450-f8e767a0] Iniciando criação de customer
13:24:50 [subscription] [CorrelationId:subscription-20250605132450-f8e767a0] Calling customer service to create customer
13:24:50 [customer] [CorrelationId:subscription-20250605132450-f8e767a0] Received create customer request
13:24:50 [customer] [customer] Customer added to in-memory collection. Email: abaaac@gmail.com
13:24:50 [customer] [customer] Saving 1 customers to database
```

### **❌ ERRADO - O que NÃO deve aparecer em {job="app"}:**
```
# Logs do Docker
2025-01-10 docker: container fc3-otl-payments.subscription-1 started

# Logs do Loki  
level=info ts=2025-01-10 msg="POST /loki/api/v1/push"

# Logs do MySQL
2025-01-10 [mysql] ready for connections

# Logs com <no value>
<no value> [subscription-1<no value>] <no value><no value>
```

## **🚨 Problemas Comuns:**

### **1. Ainda aparece <no value>:**
```bash
# Reiniciar o Promtail novamente
docker-compose restart promtail

# Aguardar 1-2 minutos para os logs aparecerem corretamente
```

### **2. Nenhum log aparece:**
```bash
# Verificar se o Promtail está rodando
docker-compose ps promtail

# Ver logs do Promtail
docker-compose logs promtail

# Reiniciar se necessário
docker-compose restart promtail
```

### **3. Logs de sistema aparecem em {job="app"}:**
- Verificar a configuração regex no `promtail/config.yaml`
- Reiniciar o Promtail

## **📊 Queries de Validação Final:**

Execute essas queries para confirmar que está tudo funcionando:

```logql
# 1. Deve retornar logs da aplicação (formato correto)
{job="app"} | limit 5

# 2. Deve retornar apenas logs do sistema  
{job="system"} | limit 5

# 3. Deve retornar apenas logs do Loki
{job="loki"} | limit 5

# 4. Esta query NÃO deve retornar nada (teste de separação)
{job="app"} |= "mysql ready for connections"

# 5. Esta query deve retornar logs de aplicação com [subscription]
{job="app"} |= "[subscription]" | limit 3

# 6. Esta query deve retornar logs com CorrelationId
{job="app"} |= "CorrelationId" | limit 5
```

## **🔍 Testando Labels Específicos:**

```logql
# Ver todos os labels disponíveis
{job="app"} | limit 1

# Filtrar por correlation_id específico (se detectado)
{job="app", correlation_id="subscription-20250605132450-f8e767a0"} | limit 10

# Filtrar por serviço da aplicação
{job="app", app_service="subscription"} | limit 10
{job="app", app_service="customer"} | limit 10
```

## **✅ Se Tudo Funcionou:**

Agora você tem:
- **Logs limpos de aplicação** com `{job="app"}` mostrando o conteúdo real dos logs
- **Logs separados por categoria** 
- **Labels extraídos automaticamente** (app_service, correlation_id)
- **Sem poluição visual** de logs de infraestrutura
- **Captura em tempo real** do Docker output

**Próximo passo:** Use as queries do `grafana-clean-log-view.md` para investigação de problemas! 