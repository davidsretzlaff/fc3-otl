
# Setup

1. Copy `.env.example` to `.env`

   ```bash
   cp .env.example .env


# Event-Driven Subscription System

This project showcases a microservices architecture built with domain events, OpenTelemetry-based observability, and a clean separation of responsibilities across services. It demonstrates realistic, loosely coupled service coordination with a focus on customer onboarding, subscription activation, and charge processing—all without a frontend.

## 🎯 Project Goals

- **Domain event modeling**: Demonstrate how to structure events that represent important business changes
- **Microservice communication via events**: Implement asynchronous and decoupled communication
- **Distributed observability**: Include tracing, logs, correlation IDs, and latency warnings
- **Realistic workflows**: Customer is created only when needed

## 🏗️ Microservices

- **Customer Service** – Manages customer creation and lookup
- **Subscription Service** – Handles subscription lifecycle
- **Charge Service** – Manages billing/charging logic

## 📡 Key Domain Events

- `SubscriptionRequested`
- `CustomerCreated`
- `CustomerVerified`
- `SubscriptionReadyForActivation`
- `ChargeSucceeded`
- `ChargeFailed`
- `SlowOperationWarning` (observability event)

## 🔄 Workflow Overview

### Step-by-step Domain Flow:

#### 1. User initiates subscription
- **Trigger**: `POST /subscriptions`
- **Service**: Subscription Service
- **Emits**: `SubscriptionRequested`

#### 2. Customer Service responds
- **Listens to**: `SubscriptionRequested`
- **Checks**: If customer exists (by email)
- **If not exists**: Creates customer
- **Emits**: `CustomerVerified` or `CustomerCreated`

#### 3. Subscription completes setup
- **Listens to**: `CustomerCreated` / `CustomerVerified`
- **Links**: Customer to subscription
- **Emits**: `SubscriptionReadyForActivation`

#### 4. Charge Service bills the user
- **Listens to**: `SubscriptionReadyForActivation`
- **Emits**: `ChargeSucceeded` or `ChargeFailed`

#### 5. Observability events
- Each service adds spans and logs
- If any operation exceeds threshold, emits: `SlowOperationWarning`

## 👁️ Observability Features

- **OpenTelemetry Collector** for centralized tracing and logging
- **Jaeger or Tempo** for trace visualization
- **Loki or Elastic** for structured JSON logs
- **Correlation IDs** passed through events and requests
- **Latency tracking** using spans
- **Custom warnings** for slow flows (`SlowOperationWarning`)

## 🛠️ Tech Stack Suggestion

| Component | Technology |
|-----------|------------|
| **Message Broker** | Kafka, NATS, or RabbitMQ |
| **Tracing** | OpenTelemetry SDK + Collector |
| **Logging** | JSON logs + FluentBit/Loki |
| **Visualization** | Grafana (for traces, logs, metrics) |
| **Runtime** | Node.js, Go, or Java |
| **Containerization** | Docker + Docker Compose |

## 🧪 How to Use / Test

Use Postman or curl to trigger:

```bash
curl -X POST http://localhost:8081/subscriptions \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "plan": "basic"}'
```

### Steps to observe:

1. **Observe** events flowing through the services
2. **Use Jaeger** to view full trace with correlation ID
3. **Review** structured logs and check for any `SlowOperationWarning` events

## 🚀 Future Improvements

- [ ] Add retry & dead-letter queue for failed events
- [ ] Add scheduled billing cycles
- [ ] Implement Saga/Process manager for full orchestration (optional)

---

**Built with focus on event-driven architecture and distributed observability** 🎯