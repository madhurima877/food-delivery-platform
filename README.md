# Real-Time Food Delivery Platform

A distributed backend system built with **Go**, demonstrating microservice architecture, event-driven communication, gRPC, Apache Kafka, Redis, PostgreSQL, fault tolerance, idempotency, distributed locking, Saga-style compensation, and Prometheus monitoring.

The project models the backend workflow of a food-delivery platform where orders move asynchronously through inventory reservation, payment processing, notifications, and driver-related services.

## Architecture

```text
                         Client
                           |
                           | HTTP
                           v
                     API Gateway
                           |
                           | gRPC
                           v
                     Order Service
                           |
                           | order.created
                           v
                         Kafka
                           |
                           v
                   Inventory Service
                           |
                           | inventory.reserved
                           v
                         Kafka
                           |
                           v
                    Payment Service
                      /          \
                     /            \
          payment.completed     payment.failed
                  |                  |
                  v                  v
        Notification/Driver    Inventory Service
                               Restore Stock
```

## Tech Stack

* **Go** — backend services
* **gRPC + Protocol Buffers** — synchronous service communication
* **Apache Kafka** — asynchronous event-driven communication
* **PostgreSQL** — persistent relational storage
* **Redis** — caching, TTL, rate limiting, and distributed locking
* **Prometheus** — metrics and monitoring
* **Docker / Docker Compose** — local infrastructure
* **Git / GitHub** — version control

## Microservices

### API Gateway

The API Gateway provides the HTTP entry point for clients and communicates with backend services using gRPC.

Responsibilities include:

* HTTP API handling
* Request routing
* gRPC client communication
* Rate limiting

Example:

```bash
curl -X POST http://localhost:8080/orders \
-H "Content-Type: application/json" \
-d '{"customer_id":"101","restaurant_id":"22","product_id":"1","quantity":1}'
```

### Order Service

Responsible for order creation and management.

Features:

* Creates and stores orders in PostgreSQL
* Exposes gRPC APIs
* Publishes `order.created` Kafka events
* Redis cache-aside implementation
* Cache TTL
* Cache invalidation after updates
* Prometheus order metrics

Current metrics include:

```text
orders_created_total
order_creation_errors_total
```

### Inventory Service

Consumes order events and manages stock reservations.

Features:

* Consumes `order.created`
* Reserves product stock transactionally
* Detects duplicate processing
* Handles insufficient stock
* Publishes `inventory.reserved`
* Kafka retry processing
* Dead Letter Queue handling
* Restores inventory after payment failure
* Prometheus monitoring

Current metrics include:

```text
inventory_reservations_total
inventory_failures_total
```

### Payment Service

Processes payments after inventory has been reserved.

Features:

* Consumes `inventory.reserved`
* PostgreSQL payment persistence
* Redis distributed locking
* Idempotent payment processing
* Publishes payment completion/failure events
* Prevents duplicate payments
* Prometheus counters and latency monitoring

Current metrics include:

```text
payments_completed_total
payments_failed_total
payment_processing_duration_seconds
```

### Notification Service

Consumes relevant events and handles notification-related processing.

### Driver Service

Handles driver-related functionality as part of the distributed order workflow.

## Kafka Event Flow

The primary successful flow is:

```text
Order Created
     |
     v
order.created
     |
     v
Inventory Service
     |
     v
inventory.reserved
     |
     v
Payment Service
     |
     v
payment.completed
```

Services communicate asynchronously through Kafka, reducing direct coupling between services.

## Kafka Reliability

The project implements reliability mechanisms for asynchronous event processing.

### Consumer Groups

Kafka consumer groups allow multiple service workers to share message processing while ensuring a partition is processed by only one consumer within the group at a time.

### Retry Mechanism

Transient failures are retried rather than immediately discarded.

Retry events are published to dedicated retry topics.

### Exponential Backoff

Retries use increasing delays to avoid continuously hammering a failing dependency.

Conceptually:

```text
Failure
  |
  +--> Retry 1
  |
  +--> Retry 2
  |
  +--> Retry 3
  |
  +--> DLQ
```

### Dead Letter Queue

After the maximum retry count is reached, failed events are sent to a Dead Letter Queue for later investigation or recovery.

Example topic:

```text
order.created.dlq
```

## Idempotency

Kafka can deliver messages more than once, so consumers must safely handle duplicate events.

### Payment Idempotency

`order_id` acts as the payment idempotency key.

PostgreSQL enforces uniqueness and payment insertion uses conflict handling:

```sql
INSERT INTO payments (...)
VALUES (...)
ON CONFLICT (order_id) DO NOTHING;
```

If the payment already exists, the repository returns a duplicate status rather than creating another payment.

This protects against repeated Kafka delivery and prevents duplicate payment records.

### Database-Level Protection

The database unique constraint provides the final idempotency guarantee even if multiple workers race to process the same event.

## Redis Distributed Lock

The Payment Service uses a Redis distributed lock before processing an order.

Conceptually:

```text
Worker A ----\
              ---> Redis Lock ---> Payment
Worker B ----/
```

Only one worker can acquire the lock for a particular order at a time.

The lock protects against concurrent processing across multiple workers or service instances.

A TTL prevents a lock from remaining forever if a worker crashes before releasing it.

### Why not only use a Go mutex?

A Go mutex protects goroutines inside a single process.

Redis is shared across service instances, allowing coordination in a distributed environment.

### Lock vs Idempotency

They solve different problems:

```text
Redis distributed lock
→ protects against concurrent processing

PostgreSQL UNIQUE constraint
→ protects against duplicate processing over time
```

## Saga-Style Compensation

The system uses event-driven compensation when payment fails.

Successful path:

```text
Order
  ↓
Reserve Inventory
  ↓
Payment
  ↓
Completed
```

Failure path:

```text
Order
  ↓
Reserve Inventory
  ↓
Payment Failed
  ↓
payment.failed
  ↓
Inventory Service
  ↓
Restore Stock
```

This prevents inventory from remaining reserved when the payment cannot be completed.

## Redis

Redis is used for several different distributed-system concerns.

### Caching

Order reads use Redis to reduce repeated database access.

```text
Request
  ↓
Check Redis
  ↓
Cache Hit ─────→ Return
  ↓ Cache Miss
PostgreSQL
  ↓
Store in Redis with TTL
  ↓
Return
```

### TTL

Cached entries automatically expire after their configured lifetime.

### Cache Invalidation

When order data changes, the corresponding cached value is removed so stale data is not returned.

### Rate Limiting

Redis is used by the API layer to control request rates.

### Distributed Locking

Redis coordinates payment processing between concurrent workers.

## Prometheus Monitoring

The project exposes Prometheus metrics through dedicated HTTP endpoints while business communication continues over gRPC.

Current development ports:

```text
Order Service Metrics      :9091/metrics
Payment Service Metrics    :9092/metrics
Inventory Service Metrics  :9093/metrics
Prometheus                 :19090
```

Prometheus currently scrapes:

```text
order-service
payment-service
inventory-service
```

### Counter Metrics

Counters measure how many times events occur.

Examples:

```text
orders_created_total
order_creation_errors_total
payments_completed_total
payments_failed_total
inventory_reservations_total
inventory_failures_total
```

### Histogram Metrics

Payment processing latency is measured using:

```text
payment_processing_duration_seconds
```

This allows payment-processing duration to be analysed using Prometheus histogram buckets.

## Graceful Shutdown

Services use Go contexts, signals, and `sync.WaitGroup` to shut down safely.

On shutdown:

```text
SIGTERM / Ctrl+C
       ↓
Cancel context
       ↓
Kafka consumers stop
       ↓
gRPC server stops gracefully
       ↓
Wait for workers
       ↓
Service exits
```

This prevents workers from being abruptly terminated during message processing.

## Testing

The project includes tests for critical reliability mechanisms.

Run all tests:

```bash
go test ./...
```

Implemented tests include:

### Redis Distributed Lock Test

Verifies that:

```text
First worker acquires lock  → true
Second acquisition          → false
Release lock
Next acquisition            → true
```

### Inventory Restoration Test

Verifies inventory is correctly restored after compensation.

### Payment Idempotency Test

The same payment is processed twice:

```text
First ProcessPayment  → COMPLETED
Second ProcessPayment → DUPLICATE
```

The database is then checked to ensure:

```text
payment rows for order = 1
```

This proves duplicate processing does not create duplicate payment records.

## Running Infrastructure

Start infrastructure:

```bash
docker-compose up -d
```

Check containers:

```bash
docker ps
```

Current local infrastructure includes:

```text
PostgreSQL  → localhost:15432
Redis       → localhost:16379
Kafka       → localhost:19092
Zookeeper   → localhost:12181
Prometheus  → localhost:19090
```

## Running Services

Run each service in a separate terminal:

```bash
go run services/order-service/main.go
go run services/inventory-service/main.go
go run services/payment-service/main.go
go run services/notification-service/main.go
go run services/driver-service/main.go
go run api-gateway/main.go
```

## Monitoring

After starting the services and Prometheus, Prometheus targets can be inspected at:

```text
localhost:19090/targets
```

Current targets should include:

```text
order-service
payment-service
inventory-service
```

Example Prometheus queries:

```text
orders_created_total
inventory_reservations_total
payments_completed_total
payments_failed_total
payment_processing_duration_seconds
```

## Project Reliability Concepts

This project currently demonstrates:

* Microservice architecture
* Event-driven architecture
* REST/HTTP API gateway
* gRPC
* Protocol Buffers
* Kafka producers and consumers
* Kafka consumer groups
* Retry handling
* Exponential backoff
* Dead Letter Queues
* Idempotent consumers
* Database unique constraints
* Race-condition protection
* Redis distributed locks
* Redis caching
* TTL
* Cache invalidation
* Rate limiting
* Saga-style compensation
* Graceful shutdown
* Concurrent Kafka workers
* Prometheus counters
* Prometheus histograms
* Service-level monitoring
* Integration/repository testing

## Remaining Roadmap

The core application and Prometheus monitoring are implemented.

Next milestones:

```text
Grafana dashboards
        ↓
OpenTelemetry distributed tracing
        ↓
Dockerize Go microservices
        ↓
Kubernetes deployment
        ↓
CI/CD pipeline
        ↓
Final integration testing
        ↓
Project cleanup and documentation
        ↓
Interview preparation
```

## Current Status

The project currently has a working end-to-end event-driven flow with:

```text
Order Creation
      ↓
Kafka
      ↓
Inventory Reservation
      ↓
Kafka
      ↓
Payment Processing
      ↓
Completion / Compensation
```

Reliability mechanisms including retries, DLQ, idempotency, distributed locking, compensation, graceful shutdown, and monitoring have been implemented and tested.

**Current project progress: approximately 66%.**
