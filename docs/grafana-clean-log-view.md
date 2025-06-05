# 📊 Como Ver Logs Limpos no Grafana - Docker Output

## **🎯 Objetivo: Logs organizados por categoria**

Agora os logs são capturados **diretamente do Docker output** (stdout/stderr) e organizados por labels:

- **`job="app"`** → Logs das aplicações (payments.subscription, payments.customer)
- **`job="system"`** → Logs do sistema (mysql, prometheus, jaeger)
- **`job="loki"`** → Logs do Loki
- **`job="container"`** → Outros containers

**Quando você filtrar por `job="app"`, verá apenas:**

```
14:54:57 [subscription] Starting CreateSubscription for joao.silva@teste.com
14:54:57 [subscription] Calling customer service to create customer  
14:54:57 [customer] [CorrelationId:subscription-20250604145457-5b4defbc] Received create customer request
14:54:57 [customer] [CorrelationId:subscription-20250604145457-5b4defbc] Database connection timeout
14:54:57 [customer] [CorrelationId:subscription-20250604145457-5b4defbc] ERROR: Failed to connect to database
14:54:57 [subscription] Customer service returned status code 500
14:54:57 [subscription] ERROR: Erro ao criar customer
```

**SEM logs do Docker, SEM logs do Loki, SEM logs do sistema!**

## **🔧 CONFIGURAÇÃO PASSO A PASSO**

### **1. Acesse o Grafana:**
- URL: `http://localhost:3000`
- Login: `admin` / `admin`

### **2. Vá para Explore:**
- Menu lateral → **Explore**
- Data source: **Loki**

### **3. Queries Organizadas por Categoria:**

#### **A) APENAS LOGS DE APLICAÇÃO (Limpos):**
```logql
{job="app"} | limit 20
```

#### **B) LOGS DE APLICAÇÃO + Filtros:**
```logql
{job="app"} |= "ERROR" | limit 10
{job="app"} |= "subscription-20250604145457-5b4defbc" | limit 30
{job="app"} |= "joao.silva@teste.com" | limit 15
```

#### **C) LOGS DE SISTEMA (se precisar debugar):**
```logql
{job="system"} | limit 20
```

#### **D) LOGS DO LOKI (para debug do próprio Loki):**
```logql
{job="loki"} | limit 20
```

#### **E) OUTROS CONTAINERS:**
```logql
{job="container"} | limit 20
```

### **4. Configuração de Visualização:**

#### **A) No painel Options (lado direito):**
- **Max data points**: `50`
- **Order**: `Time (ascending)`

#### **B) Na aba Log (abaixo da query):**
- **Display mode**: `Logs`
- **Wrap lines**: `ON`
- **Show time**: `ON`
- **Show labels**: `OFF` (menos poluição visual)

## **💡 QUERIES ESPECÍFICAS PARA APLICAÇÃO:**

### **1. Investigação de Erro por Email:**
```logql
{job="app"} |= "joao.silva@teste.com" |= "ERROR" | limit 10
```

### **2. Fluxo Completo por Correlation ID:**
```logql
{job="app"} |= "subscription-20250604145457-5b4defbc" | limit 50
```

### **3. Apenas Erros da Aplicação:**
```logql
{job="app"} |= "ERROR" | limit 20
```

### **4. Por Serviço Específico:**
```logql
{job="app", service="subscription"} | limit 20
{job="app", service="customer"} | limit 20
```

### **5. Por Nível de Log:**
```logql
{job="app", level="error"} | limit 15
{job="app", level="info"} | limit 30
```

## **🚀 WORKFLOW PRÁTICO DE INVESTIGAÇÃO:**

### **Passo 1: Query Inicial (só aplicação)**
```logql
{job="app"} |= "joao.silva@teste.com" | limit 10
```

### **Passo 2: Extrair Correlation ID**
Do resultado, copie o `correlation_id`: `subscription-20250604145457-5b4defbc`

### **Passo 3: Ver Fluxo Completo**
```logql
{job="app"} |= "subscription-20250604145457-5b4defbc" | limit 50
```

### **Passo 4: Focar nos Erros**
```logql
{job="app"} |= "subscription-20250604145457-5b4defbc" |= "ERROR" | limit 10
```

## **📋 LABELS DISPONÍVEIS:**

### **Para job="app" você tem:**
- **`service`**: `subscription`, `customer`
- **`level`**: `info`, `error`, `warning`, `debug`
- **`correlation_id`**: ID único da operação
- **`operation`**: nome da operação sendo executada

### **Exemplos de uso dos labels:**
```logql
# Apenas erros do serviço de subscription
{job="app", service="subscription", level="error"}

# Todas as operações de CreateSubscription
{job="app", operation="CreateSubscription"}

# Logs de um correlation_id específico
{job="app", correlation_id="subscription-20250604145457-5b4defbc"}
```

## **🎯 QUERIES AVANÇADAS:**

### **1. Contar Erros por Serviço:**
```logql
sum(count_over_time({job="app", level="error"}[5m])) by (service)
```

### **2. Taxa de Erro da Aplicação:**
```logql
sum(rate({job="app", level="error"}[5m])) / sum(rate({job="app"}[5m])) * 100
```

### **3. Operações Mais Lentas:**
```logql
{job="app"} |= "duration_ms" | json | duration_ms > 3000
```

## **⚡ VANTAGENS DA NOVA CONFIGURAÇÃO:**

### **✅ Logs Limpos:**
- **SEM logs do Docker** (`docker: container started`)
- **SEM logs do Loki** (`level=info msg="..."`)
- **SEM logs do sistema** (`mysql connection pooling`)
- **APENAS logs da sua aplicação!**

### **✅ Organização por Labels:**
- Fácil filtragem por categoria
- Queries mais específicas e rápidas
- Menos ruído visual

### **✅ Captura em Tempo Real:**
- Logs vêm direto do `stdout/stderr` dos containers
- Não depende de arquivos de log
- Funciona perfeitamente com containers

## **🔍 RESULTADO ESPERADO:**

**Ao usar `{job="app"}`, você verá apenas:**

```
14:54:57 [subscription] Starting CreateSubscription for joao.silva@teste.com
14:54:57 [subscription] Calling customer service to create customer
14:54:57 [customer] [CorrelationId:subscription-20250604145457-5b4defbc] Received create customer request
14:54:57 [customer] [CorrelationId:subscription-20250604145457-5b4defbc] Database connection timeout
14:54:57 [customer] [CorrelationId:subscription-20250604145457-5b4defbc] ERROR: Failed to connect to database
14:54:57 [subscription] Customer service returned status code 500
14:54:57 [subscription] ERROR: Erro ao criar customer
```

**Nenhum log de infraestrutura ou sistema!**

## **🛠️ Para Aplicar as Mudanças:**

1. **Reinicie o Promtail:**
```bash
docker-compose restart promtail
```

2. **Aguarde 30 segundos** para os logs começarem a aparecer

3. **Teste no Grafana:**
```logql
{job="app"} | limit 10
```

---

**💡 DICA FINAL:** Agora você tem logs organizados e limpos! Use `{job="app"}` para investigação de aplicação e `{job="system"}` apenas quando precisar debugar a infraestrutura. 