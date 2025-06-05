# 🔍 Workflow de Investigação de Logs - Do Erro ao Root Cause

## 🚨 Cenário: "Usuário tentou criar subscription e deu erro"

### **PASSO 1: Começar pela informação que você TEM**

#### A) Cliente reportou erro → Busque pelo email
```logql
{job="containers"} |= "email_do_cliente@exemplo.com"
```

#### B) Erro aconteceu "agora há pouco" → Busque erros recentes
```logql
{job="containers"} |= "ERROR"
```

#### C) Erro em operação específica → Busque pela operação
```logql
{job="containers"} |= "CreateSubscription"
{job="containers"} |= "ActivateSubscription"
{job="containers"} |= "subscription" |= "ERROR"
```

#### D) Erro em horário específico → Ajuste o time range
```
Time Range: Last 30 minutes (ou quando o erro aconteceu)
```

---

### **PASSO 2: ENCONTRAR o Correlation ID**

#### Query Combinada - Erro de Subscription Recente:
```logql
{job="containers"} |= "subscription" |= "ERROR"
```

**No resultado, você vai ver algo como:**
```json
{
  "timestamp": "2025-01-04T02:30:15Z",
  "level": "error", 
  "service": "subscription",
  "correlation_id": "subscription-20250604023015-abc123",
  "customer_email": "joao@exemplo.com",
  "operation": "activate_subscription",
  "error": "customer service returned status code 500"
}
```

**🎯 PEGUE o `correlation_id`: `subscription-20250604023015-abc123`**

---

### **PASSO 3: RASTREAR TODO O FLUXO**

Agora que você tem o correlation ID, rastreie TUDO:

```logql
{job="containers"} |= "subscription-20250604023015-abc123"
```

**Você vai ver a sequência completa:**
```
14:30:10 [subscription] Starting activate_subscription for joao@exemplo.com
14:30:11 [subscription] Calling customer service to create customer
14:30:12 [customer] [CorrelationId:subscription-20250604023015-abc123] Received create customer request
14:30:13 [customer] [CorrelationId:subscription-20250604023015-abc123] Database connection failed
14:30:13 [customer] [CorrelationId:subscription-20250604023015-abc123] Error: connection timeout
14:30:14 [subscription] Customer service returned status code 500
14:30:14 [subscription] ERROR: Failed to activate subscription
```

**🔍 ROOT CAUSE ENCONTRADO: Database connection failed no Customer Service**

---

## 🎯 **Queries por Cenário de Investigação**

### **📧 Cenário 1: "Cliente X não conseguiu criar subscription"**
```logql
# Passo 1: Buscar por cliente específico
{job="containers"} |= "joao@exemplo.com"

# Passo 2: Focar nos erros desse cliente
{job="containers"} |= "joao@exemplo.com" |= "ERROR"

# Passo 3: Ver o fluxo completo (pegue o correlation_id do resultado acima)
{job="containers"} |= "correlation_id_encontrado"
```

### **⏰ Cenário 2: "Sistema lento às 14h30"**
```logql
# Passo 1: Operações lentas nesse horário
{job="containers"} |= "duration_ms" | json | duration_ms > 3000

# Passo 2: Ver o que estava acontecendo
{job="containers"} |= "ERROR"  # No time range de 14h25-14h35

# Passo 3: Investigar correlation_ids dos problemas encontrados
```

### **🚨 Cenário 3: "Muitos erros 500 agora"**
```logql
# Passo 1: Ver todos os 500s recentes
{job="containers"} |= "status code 500"

# Passo 2: Agrupar por serviço para ver onde está o problema
{job="containers"} |= "status code 500" |= "customer"  # ou "subscription"

# Passo 3: Pegar correlation_ids dos erros e investigar cada fluxo
```

### **🔄 Cenário 4: "Integração entre serviços falhando"**
```logql
# Passo 1: Erros de comunicação
{job="containers"} |= "customer_service_error"
{job="containers"} |= "connection"
{job="containers"} |= "timeout"

# Passo 2: Ver ambos os lados da comunicação
{job="containers"} |= "correlation_id_encontrado"
```

---

## 🛠️ **Workflow Completo de Investigação**

### **1. 🎯 IDENTIFICAR o Problema**
- [ ] Que operação falhou?
- [ ] Qual cliente foi afetado?
- [ ] Quando aconteceu?
- [ ] Que erro foi reportado?

### **2. 🔍 BUSCAR os Logs Iniciais**
```logql
# Use uma dessas estratégias:
{job="containers"} |= "email_do_cliente"           # Se tem o cliente
{job="containers"} |= "ERROR"                      # Se quer ver erros gerais  
{job="containers"} |= "operation_name" |= "ERROR" # Se sabe a operação
{job="containers"} |= "status code 500"           # Se tem código HTTP
```

### **3. 📋 EXTRAIR o Correlation ID**
- No resultado da query acima, procure por:
  - `"correlation_id": "subscription-..."`
  - `[CorrelationId:subscription-...]`
- Copie o correlation ID completo

### **4. 🌊 RASTREAR o Fluxo Completo**
```logql
{job="containers"} |= "SEU_CORRELATION_ID_AQUI"
```

### **5. 🔬 ANALISAR a Sequência**
- Ordene por timestamp
- Identifique onde o erro começou
- Veja que serviço falhou primeiro
- Verifique os códigos de status
- Analise os tempos de resposta

### **6. 🎯 IDENTIFICAR Root Cause**
- Último log de sucesso
- Primeiro log de erro  
- Serviço que originou o problema
- Erro específico (connection, timeout, validation, etc.)

---

## 📝 **Exemplo Prático Completo**

**Situação:** "Cliente joao@exemplo.com tentou ativar subscription e deu erro 500"

**1. Query inicial:**
```logql
{job="containers"} |= "joao@exemplo.com" |= "ERROR"
```

**2. Resultado encontrado:**
```
{"correlation_id":"subscription-20250604141530-xyz789", "customer_email":"joao@exemplo.com", "error":"status code 500"}
```

**3. Query de investigação:**
```logql
{job="containers"} |= "subscription-20250604141530-xyz789"
```

**4. Análise do fluxo:**
```
14:15:30 [subscription] Starting activation...
14:15:31 [subscription] Calling customer service... 
14:15:35 [customer] Database timeout after 3 seconds
14:15:35 [customer] Error: failed to create customer
14:15:36 [subscription] Error: customer service returned 500
```

**5. Root Cause:** Database timeout no Customer Service

**6. Ação:** Investigar performance/conectividade do banco de dados

---

## 🚀 **Dicas para Investigação Eficiente**

### ✅ **FAÇA:**
- Comece sempre pela informação que você TEM (cliente, horário, operação)
- Use time ranges apropriados (últimos 30min, última hora)
- Combine queries: `|= "ERROR" |= "subscription"`
- Sempre pegue o correlation ID para rastrear o fluxo completo

### ❌ **NÃO FAÇA:**
- Não comece com correlation ID (você não vai ter)
- Não ignore o contexto temporal
- Não foque apenas em um serviço
- Não esqueça de verificar ambos os lados da comunicação

### 🎯 **LEMBRE-SE:**
O correlation ID é o **resultado** da investigação inicial, não o **ponto de partida**! 