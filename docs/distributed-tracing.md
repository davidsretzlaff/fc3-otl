# Distributed Tracing and Context Propagation

## Basic Concepts

### Trace
A `trace` represents the complete journey of a request through a distributed system. Think of it as the "path" that a request takes through different services.

### Span
A `span` represents a single operation or unit of work within a trace. For example:
- An HTTP call
- A database query
- A specific processing task

### Important IDs

#### TraceID
- A unique identifier for the entire trace
- Remains the same across all services the request passes through
- Allows viewing the complete journey of the request
- Example: `47a4d9efca74d6b962c64bd5a0d4f83d`

#### SpanID
- Uniquely identifies a specific span
- Each operation within the trace has its own SpanID
- Example: `2d281dd73dcf7574`

#### ParentSpanID
- References the SpanID of the span that originated the current operation
- Allows building the hierarchy/tree of spans
- Example: `c42a0b27839aa914`

## Context Propagation

### How It Works
1. **Service A (Origin)**
   ```go
   // Go (using OpenTelemetry)
   ctx, span := tracer.Start(ctx, "OperationA")
   defer span.End()
   
   // Inject context into HTTP headers
   propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
   ```

2. **HTTP Headers**
   ```
   traceparent: 00-47a4d9efca74d6b962c64bd5a0d4f83d-c42a0b27839aa914-01
   ```
   Format: `version-traceID-spanID-flags`

3. **Service B (Destination)**
   ```csharp
   // C# (using OpenTelemetry)
   // ASP.NET Core automatically extracts context from headers
   app.UseOpenTelemetry();
   ```

### Propagation Headers

#### traceparent
- Carries essential trace information
- Format: `version-traceID-spanID-flags`
- Example: `00-47a4d9efca74d6b962c64bd5a0d4f83d-c42a0b27839aa914-01`

#### tracestate
- Carries additional trace information
- Used for vendor-specific data
- Optional

## CorrelationID vs TraceID

### CorrelationID
- Unique identifier to correlate logs and events
- Generally simpler, just an ID
- Frequently used in structured logging
- Does not carry hierarchy information

### TraceID
- Part of a more complex tracing system
- Carries hierarchy information (ParentSpanID)
- Allows detailed flow visualization
- Integrated with APM (Application Performance Monitoring) tools

### When to Use Each?
- **CorrelationID**: For simpler systems where you just need to correlate logs
- **TraceID**: For complex distributed systems where you need to understand the complete flow

## Practical Example

### Subscription -> Customer Flow

1. **Subscription Initiates the Trace**
   ```go
   ctx, span := tracer.Start(ctx, "CreateSubscription")
   ```

2. **Subscription Calls Customer**
   ```go
   // Creates a new child span
   ctx, span := tracer.Start(ctx, "CustomerClient.CreateCustomer")
   
   // Headers are automatically injected
   propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
   ```

3. **Customer Receives the Request**
   ```csharp
   // ASP.NET Core automatically extracts the context
   // The new span will have:
   // - Same TraceID as Subscription
   // - Its own new SpanID
   // - ParentSpanID referencing the Subscription span
   ```

## Visualization in Jaeger

```
Trace
├── Subscription: /subscriptions (SpanID: abc123)
│   └── CustomerClient.CreateCustomer (SpanID: def456)
│       └── Customer: POST /api/customer (SpanID: ghi789)
│           └── Customer: Processing (SpanID: jkl012)
```

## Benefits

1. **Distributed Debugging**
   - Visualize the complete request flow
   - Identify bottlenecks and failures
   - Understand service dependencies

2. **Performance**
   - Measure response times
   - Identify slow operations
   - Optimize based on real data

3. **Observability**
   - Monitor distributed systems
   - Understand production behavior
   - Facilitate problem resolution

## Popular Tools

1. **OpenTelemetry**
   - Open-source instrumentation standard
   - Multi-language support
   - Exporters for various backends

2. **Jaeger**
   - Visual interface for traces
   - Performance analysis
   - Advanced filters and searches

3. **Zipkin**
   - Alternative to Jaeger
   - Focus on low latency
   - Good Spring integration

## Debugging and Troubleshooting

### Verifying Headers and Context

1. **Logs in Origin Service (Go)**
   ```go
   // Add logs before sending the request
   fmt.Printf("Headers being sent: %+v\n", req.Header)
   ```

2. **Logs in Destination Service (C#)**
   ```csharp
   app.Use(async (context, next) =>
   {
       var activity = System.Diagnostics.Activity.Current;
       if (activity != null)
       {
           // Log tracing headers
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

### Docker Debug Commands

1. **Check logs of a specific service**
   ```bash
   # View customer service logs
   docker-compose logs payments.customer

   # View subscription service logs
   docker-compose logs payments.subscription

   # View logs in real-time
   docker-compose logs -f payments.customer
   ```

2. **Check logs inside the container**
   ```bash
   # Access logs inside the container
   docker exec fc3-otl-payments.customer-1 cat /app/logs/app.log
   ```

3. **Rebuild containers after changes**
   ```bash
   # Rebuild a specific service
   docker-compose up -d --build payments.customer
   docker-compose up -d --build payments.subscription

   # Rebuild all services
   docker-compose up -d --build
   ```

### What to Look for in Logs

1. **In Origin Service**
   - `traceparent` headers being injected
   - Generated TraceID
   - Current span's SpanID

2. **In Destination Service**
   - Received `traceparent` headers
   - Same TraceID as origin service
   - ParentSpanID matching origin span's SpanID
   - New generated SpanID

### Common Issues

1. **Headers Not Arriving**
   - Check if propagator is correctly configured
   - Confirm header injection is happening
   - Verify no proxy/gateway is removing headers

2. **Different TraceID**
   - Check if context is being properly propagated
   - Confirm no multiple OpenTelemetry configurations

3. **Spans Not Appearing in Jaeger**
   - Check OpenTelemetry collector configuration
   - Confirm spans are being ended (defer span.End())
   - Verify exporter configuration

### Debugging Tips

1. **Use Span Tags**
   ```go
   span.SetAttributes(
       attribute.String("customer.name", request.Name),
       attribute.String("customer.email", request.Email),
   )
   ```

2. **Add Span Events**
   ```go
   span.AddEvent("starting_http_call")
   // ... code ...
   span.AddEvent("http_call_completed")
   ```

3. **Capture Errors**
   ```go
   if err != nil {
       span.RecordError(err)
       span.SetStatus(codes.Error, err.Error())
   }
   ```

## Event-Driven Tracing

### Overview
In event-driven architectures, context propagation works differently from HTTP requests. Instead of HTTP headers, the trace context needs to be included in the event/message payload or metadata.

### Common Message Brokers Support
1. **Kafka**
   - Uses message headers to carry trace context
   - OpenTelemetry provides built-in propagators

2. **RabbitMQ**
   - Uses message properties to carry trace context
   - Can use AMQP headers or message properties

3. **AWS SQS/SNS**
   - Uses message attributes to carry trace context
   - Supports message system attributes

### Implementation Examples

1. **Publishing Events (Go with Kafka)**
   ```go
   func PublishEvent(ctx context.Context, event MyEvent) error {
       // Create a span for the publish operation
       ctx, span := tracer.Start(ctx, "PublishEvent")
       defer span.End()

       // Create message headers
       headers := []kafka.Header{}
       
       // Inject trace context into Kafka headers
       otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))

       // Create and send the message
       msg := kafka.Message{
           Headers: headers,
           Key:     []byte(event.ID),
           Value:   eventData,
       }
       return producer.Produce(&msg)
   }
   ```

2. **Consuming Events (Go with Kafka)**
   ```go
   func ConsumeEvent(msg *kafka.Message) error {
       // Extract trace context from Kafka headers
       headers := propagation.HeaderCarrier(msg.Headers)
       ctx := otel.GetTextMapPropagator().Extract(context.Background(), headers)

       // Create a span for the consume operation
       ctx, span := tracer.Start(ctx, "ConsumeEvent")
       defer span.End()

       // Process the event...
       return processEvent(ctx, msg.Value)
   }
   ```

3. **RabbitMQ Example (C#)**
   ```csharp
   // Publishing
   public async Task PublishEvent(MyEvent evt)
   {
       using var activity = _activitySource.StartActivity("PublishEvent");
       
       var props = channel.CreateBasicProperties();
       props.Headers = new Dictionary<string, object>();
       
       // Inject trace context into message headers
       var propagator = Propagators.DefaultTextMapPropagator;
       propagator.Inject(new PropagationContext(activity.Context, Baggage.Current),
           props.Headers,
           (headers, key, value) => headers[key] = value);

       await channel.BasicPublishAsync(exchange, routingKey, props, eventData);
   }

   // Consuming
   public async Task HandleMessage(IModel channel, BasicDeliverEventArgs ea)
   {
       var headers = ea.BasicProperties.Headers;
       
       // Extract trace context from message headers
       var parentContext = Propagators.DefaultTextMapPropagator
           .Extract(default, headers, (headers, key) => 
               headers.TryGetValue(key, out var value) ? value.ToString() : null);

       using var activity = _activitySource.StartActivity(
           "ConsumeEvent", 
           ActivityKind.Consumer,
           parentContext.ActivityContext);

       // Process the message...
   }
   ```

### Best Practices for Event-Driven Tracing

1. **Message Enrichment**
   - Include trace context in message metadata/headers
   - Avoid modifying the actual message payload
   - Use standardized header names (e.g., `traceparent`)

2. **Correlation Strategy**
   - Use event/message ID as correlation ID
   - Link related spans using span links
   - Maintain causality chain across events

3. **Handling Batch Processing**
   ```go
   func ProcessBatch(messages []Message) {
       for _, msg := range messages {
           // Extract context from message
           ctx := extractContext(msg)
           
           // Create a new span but link it to the message context
           newCtx, span := tracer.Start(context.Background(), "ProcessMessage",
               trace.WithLinks(trace.Link{Context: ctx}))
           
           // Process message...
           span.End()
       }
   }
   ```

4. **Dead Letter Queues (DLQ)**
   ```go
   func HandleFailedMessage(msg Message, err error) {
       // Extract original context
       originalCtx := extractContext(msg)
       
       // Create new span for DLQ operation but link to original
       ctx, span := tracer.Start(context.Background(), "MoveToDLQ",
           trace.WithLinks(trace.Link{Context: originalCtx}))
       defer span.End()

       span.SetAttributes(
           attribute.String("error", err.Error()),
           attribute.String("original_queue", msg.Queue),
           attribute.String("message_id", msg.ID),
       )

       // Move to DLQ...
   }
   ```

### Common Challenges in Event-Driven Tracing

1. **Asynchronous Nature**
   - Events may be processed out of order
   - Multiple consumers may process the same event
   - Need to handle concurrent processing

2. **Message Transformation**
   - Preserve trace context during message transformation
   - Handle protocol conversions (e.g., Kafka to RabbitMQ)
   - Maintain context across different message formats

3. **Retry Mechanisms**
   - Preserve original trace context during retries
   - Link retry spans to original processing span
   - Track retry count and delays in spans

4. **Monitoring and Debugging**
   - Use span attributes to track queue metrics
   - Monitor processing delays and backlogs
   - Track message lifecycle across systems 