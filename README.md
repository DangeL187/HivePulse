# About

**HivePulse** is a high-load, event-driven system designed for real-time telemetry collection from remote devices.
Built as a full-fledged microservices ecosystem, it efficiently ingests, processes, and stores massive streams of data
while maintaining high availability and horizontal scalability.

**HivePulse** reliably handles around **50K RPS**, demonstrates predictable behavior under heavy load, and is built for
robust, fault-tolerant operation.

# 🔥 Services Overview

## 🌐 MQTT Broker

- EMQX as an MQTT Broker

## ⚡️ Ingress (MQTT-to-Kafka) service

The **ingress** service is responsible for collecting telemetry data from devices via **MQTT** and forwarding it to **Kafka** for downstream processing.
It is built with modularity in mind, allowing easy replacement or extension of components without affecting the rest of
the system.

### Architecture Overview

1. **ConsumerLoop**
    - Reads messages from the configured `consumer` module (e.g., `MQTTConsumer`).
    - Publishes incoming messages into a shared channel `msgChan` (buffer size 10,000) for further processing.
2. **ProducerLoop**
    - Runs multiple worker goroutines reading from `msgChan`.
    - Processes messages using `ProducerLoop.processMessage`.
    - Authenticates devices via the `AuthService` module.
    - Upon successful authentication, sends messages to the configured `producer` module (e.g., `KafkaProducer`).
    - One dedicated worker listens to the `producer` module's error channel for monitoring and retries.
3. **AuthService**
    - Authenticates devices using their JWT tokens via the authenticator module (e.g., `GRPCAuthenticator`).
    - Runs a background worker that reads authentication error events and notifies devices through the `publisher`
      module (e.g., `MQTTPublisher`).
    - Offloads token validation from the **auth** service by using public JWT tokens, reducing the load on the central
      auth system. A caching mechanism is planned to further optimize token validation per message.
4. **GRPCAuthenticator**
    - Uses the public JWT token obtained from the Auth service for device authentication.
    - Reduces repetitive calls to the **auth** service and allows the **ingress** service to scale independently.

### Key Features

- **Modular Design**: Components do not depend directly on each other, enabling easy swapping or extension of
  implementations.
- **High Throughput**: Efficient message channeling and worker pooling support large-scale telemetry ingestion.
- **Fault-Tolerant Authentication**: Background error handling ensures devices are notified promptly about
  authentication issues.
- **Monitoring**: Prometheus metrics endpoint for observability and performance tracking.

## ⚡️ Consumer (Kafka-to-ClickHouse) service

The **consumer** service is a horizontally scalable service responsible for reading telemetry events from **Kafka** and
flushing them into **ClickHouse** in batches.
Multiple instances of this service can run simultaneously, each consuming a share of **Kafka** partitions for balanced
load distribution.

### Architecture Overview

1. **ConsumerLoop**
    - Reads messages from the configured `consumer` module (e.g., `KafkaConsumer`).
    - Publishes incoming messages into a shared channel `msgChan` (buffer size 10,000) for batch processing.
2. **MessageBatchFlusher**
    - Runs multiple worker goroutines that read from `msgChan`.
    - Aggregates messages into batches.
    - Writes batches to `ClickHouse` via the configured `flusher` module (e.g., `KafkaClickHouseFlusher`).
3. **KafkaConsumer**
    - Marks messages as read every second, reducing latency compared to acknowledging each message individually.
    - `Kafka` topic is created with 12 partitions, allowing even load distribution across multiple **consumer** service
      instances.

### Key Features

- **Modular Design**: Components are loosely coupled, enabling easy replacement or extension of consumer and flusher
  implementations.
- **Replicable and Scalable**: Multiple service instances can consume different `Kafka` partitions in parallel.
- **Batch Processing**: Efficiently flushes messages in bulk to `ClickHouse`, minimizing write overhead.
- **Monitoring**: Prometheus metrics endpoint for observability and performance tracking.

## 🔒️ Auth Service

The **auth** service is responsible for authenticating and authorizing both users and devices.
It exposes **REST** endpoints for user management and **gRPC** endpoints for device authentication, providing a secure
gateway for HivePulse's telemetry ecosystem.

### Architecture Overview

1. **Users and Devices Authentication and Authorization** (REST, port 8000)
    - Users and devices are stored in **PostgreSQL**.
    - Authentication uses **JWT** tokens (ed25519) with public/private keys, supporting access and refresh tokens.
    - User roles are managed via **Casbin**:
        - `admin`: can grant/revoke roles for users
        - `operator`: can register devices and monitor them
    - **Endpoints**:
        - `POST /users/login` - user login
        - `POST /users/register` - user registration
        - `GET /users/:id/roles` - get roles for a user
        - `POST /users/:id/roles` - grant role
        - `DELETE /users/:id/roles/:role` - revoke role
        - `POST /devices/login` - device login
        - `POST /devices/register` - device registration
        - `POST /devices/refresh` - refresh device token
    - **How to** generate keys:
      ```bash
      openssl genpkey -algorithm Ed25519 -out private.pem
      openssl pkey -in private.pem -pubout -out public.pem
      
      openssl pkey -in private.pem -outform DER | base64 -w0
      openssl pkey -in public.pem -pubin -outform DER | base64 -w0
      ```
2. **Testing**:
    - Lightweight **Python** tests are included for basic functionality.
    - **Important**: run tests before starting the service to ensure DB and key setup is correct.

### Key Features

- **Authentication**: Securely handles both users and devices.
- **Role-Based Access Control**: **Casbin** ensures fine-grained authorization.
- **Secure Tokens**: **JWT** with ed25519 keys for high security.
- **Ready for Integration**: Provides **gRPC** interface for other services (e.g., **ingress**) to validate devices
  efficiently.

## 📈 Metrics

You can find simple **Grafana** dashboard for **ClickHouse** [here](https://github.com/DangeL187/HivePulse/blob/main/metrics/clickhouse.json).

# 🔧 Code Architecture Decisions

## Architectural Approach

- **Feature-Sliced Design (FSD)**
    - Chosen deliberately for its effectiveness in large codebases.
    - Works seamlessly with **Clean Architecture** principles, ensuring modularity, separation of concerns, and easy
      extensibility.
- **Device Service as Load Simulator**
    - Initially was meant to be a one-file mock program for sending metrics.
    - The code does not follow strict architectural patterns and should be treated as a simulation tool.
    - By default, 200 replicas of the **device** service run, each sending 100 metrics/sec, generating roughly 20K
      requests/sec for testing system load.

## Libraries and Tooling
- **Error Handling**: Custom library [erax](https://github.com/DangeL187/erax) is used for consistent and convenient error handling across all services.
- **Logging**: Uber Zap is used for structured and performant logging.
