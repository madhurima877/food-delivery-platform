API Gateway :8080
    │
    ├── gRPC → Order Service        :50051 ✅
    ├── gRPC → Inventory Service    :50052 ✅
    ├── gRPC → Payment Service      :50053 ✅
    ├── gRPC → Notification Service :50054 ✅
    └── gRPC → Driver Service       :50056 ✅

    use context
    graceful shutdown

    Transactional Outbox Pattern