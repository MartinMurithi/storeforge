# StoreForge Technical Specification Document

## 1. Introduction
StoreForge is a microservices-based Platform-as-a-Service (PaaS) enabling merchants to launch fully functional e-commerce stores in under 5 minutes. The platform is designed for scalability, modularity, and operational reliability, supporting a trial-first subscription model and direct merchant payouts via payment providers like M-Pesa and Stripe.

This document serves as the reference architecture for development teams and outlines the entire technical ecosystem, including services, responsibilities, communication patterns, infrastructure, and production considerations.

---

## 2. System Overview

### 2.1 Architectural Principles
- **Domain-Driven Design (DDD):** Each bounded context owns its data and business logic.
- **Microservices Isolation:** Every service is independently deployable and scalable.
- **Event-Driven Workflows:** Asynchronous events orchestrate business processes.
- **Gateway Edge:** All client interactions pass through a secure gateway.
- **Zero-Touch Onboarding:** Merchants can start selling with minimal technical involvement.
- **Plugin/Extension Ecosystem:** Optional features integrate via events without modifying core services.

### 2.2 Communication Patterns
- **Synchronous:** gRPC between microservices for real-time queries.
- **Asynchronous:** Message broker (Kafka, NATS, RabbitMQ) for events like `StoreProvisioned`, `OrderCompleted`.
- **Client-facing:** REST or GraphQL via a custom Go gateway.

---

## 3. Infrastructure Components

### 3.1 NGINX (Edge Reverse Proxy)
**Responsibilities:**
- SSL/TLS termination
- URL-based routing to gateway/project modules
- Load balancing across multiple gateway instances
- Rate limiting and security enforcement
- Optional caching of static or frequent responses

**Key Principle:** Keeps internal microservices hidden from external exposure.

### 3.2 API Gateway
**Implemented in:** Go (Gin or Chi framework)

**Responsibilities:**
- Single entry point for HTTP/REST or GraphQL requests
- JWT token validation and role-based access enforcement
- Aggregates data from multiple services for dashboards
- Translates HTTP requests into internal gRPC calls
- Publishes asynchronous events for workflow orchestration
- Handles cross-cutting concerns: logging, rate-limiting, circuit breaking, request validation
- Optional real-time communication via WebSocket or SSE

**Notes:**
- Modular design to support multiple projects if needed.
- Optional feature flags can be enforced at this level.

---

## 4. Core Domain Services

### 4.1 Onboarding Service
**Purpose:** Manages merchant lifecycle and store provisioning.

**Responsibilities:**
- Merchant registration confirmation
- Store setup: URL, theme, default products
- Trial period management
- Orchestrates provisioning workflow via asynchronous events
- Publishes: `MerchantVerified`, `StoreProvisioned`

**Communication:** Receives gRPC requests from gateway; publishes events to supporting domains.

### 4.2 Auth Service
**Purpose:** Identity, authentication, and authorization.

**Responsibilities:**
- User registration, login, password hashing
- OTP/email verification
- Role and permission management
- JWT token issuance

**Communication:**
- Gateway verifies JWTs issued by Auth Service
- Services can query Auth for role verification if necessary

---

## 5. Supporting Domain Services

### 5.1 Catalog Service
- Manages products, categories, and inventory
- Handles default catalog creation during onboarding
- Publishes inventory change events for analytics or notifications

### 5.2 Payments Service
**Purpose:** Handles all merchant payment setup, processing, and integration with external payment providers.

**Responsibilities:**

1. **Payment Method Integration**
    - Supports M-Pesa, Stripe, and other future plugins.
    - Differentiates sandbox mode (trial merchants) vs live mode (paid subscriptions).

2. **Merchant Payment Setup**
    - **M-Pesa:** Merchants must register with Safaricom to obtain a Paybill or Till number. This number links the merchant’s store to their account, ensuring funds go directly to them. Merchants configure this in the dashboard; credentials are securely stored.
    - **Stripe:** Merchants register a Stripe account and provide API keys. Platform provides sandbox/test accounts for trial merchants.

3. **Transaction Processing**
    - Processes orders and routes funds directly to merchant accounts.
    - Publishes events: `PaymentCompleted`, `PaymentFailed`, `PaymentRefunded`.

4. **Sandbox vs Live Flow**
    - Trial Merchants: Platform-assigned sandbox credentials, no real funds processed.
    - Paid Merchants: Must provide live credentials (Paybill/Till or Stripe API keys) for real transactions.

5. **Security & Compliance**
    - Payment credentials encrypted and stored securely.
    - Trial and live flows separated to prevent accidental fund access.
    - Platform does not handle merchant funds directly.

**Integration Points:**
- Gateway for configuration/testing
- Onboarding Service triggers initial payment setup
- Orders Service validates payment completion
- Notifications Service alerts merchants of payment events
- Analytics Service logs payment data for dashboards

### 5.3 Storefront / Provisioning Service
- Handles store URL assignment, theme application, hosting configuration
- Triggers asynchronous events for store readiness and analytics

### 5.4 Notifications Service
- Sends emails, SMS, OTPs, and webhooks
- Reacts to events like `PaymentCompleted`, `OrderShipped`, `StoreProvisioned`

### 5.5 Analytics / Observability Service
- Aggregates metrics, logs, and performance data
- Provides dashboard data for merchants
- Subscribes to events across all services

### 5.6 Subscription / Feature Flags Service
- Manages trial periods, plan upgrades, and feature access
- Publishes `SubscriptionUpgraded` events for dependent services

### 5.7 Extensions / Plugin Registry Service
- Optional integrations like Fleetbase delivery, advanced analytics, marketing tools
- Plugins react asynchronously to events from core services
- Allows monetization of add-on features without modifying core services

---

## 6. Operational / Supporting Entities

### 6.1 Orders Service
- Transactional hub for merchant sales
- Responsibilities:
  - Validate product availability via Catalog Service
  - Process payments via Payments Service
  - Publish events: `OrderCreated`, `OrderPaid`, `OrderShipped`, `OrderCompleted`
- Database: relational (PostgreSQL) for ACID compliance

### 6.2 Customers Service
- Manages buyer profiles and purchase history
- Publishes events for analytics and marketing
- Optional support for segmentation, recommendations, and loyalty programs

---

## 7. Infrastructure & Operational Services

### 7.1 Event Broker
- Kafka, NATS, or RabbitMQ
- Handles asynchronous workflows and plugin integration

### 7.2 File / Media Storage
- Stores product images, theme assets, logos
- Typically S3-compatible storage + CDN for delivery

### 7.3 Background Jobs / Scheduler
- Processes delayed or recurring tasks:
  - Trial expirations
  - Reports generation
  - Retry failed payments

### 7.4 Audit & Compliance
- Maintains logs of critical actions: login, orders, payment events, subscription changes
- Ensures traceability and regulatory compliance

### 7.5 Search / Discovery Service
- Handles product, theme, and plugin search for merchants
- Optional ElasticSearch or Postgres full-text search

### 7.6 Security & Secrets Management
- Centralized store for API keys (Stripe, M-Pesa)
- Rotates credentials, ensures encrypted storage
- Prevents accidental exposure of merchant credentials

---

## 8. Databases & Storage Strategy

| Service    | Database Type         | Notes                                    |
| ---------- | --------------------- | ---------------------------------------- |
| Auth       | PostgreSQL            | User credentials, roles, JWT secrets     |
| Onboarding | PostgreSQL            | Merchant and store lifecycle             |
| Catalog    | PostgreSQL/MongoDB    | Products, categories, inventory          |
| Payments   | PostgreSQL            | Transaction ledger, payment status       |
| Orders     | PostgreSQL            | Order lifecycle, transactional integrity |
| Customers  | PostgreSQL/MongoDB    | Buyer data, purchase history             |
| Analytics  | Clickhouse/PostgreSQL | Aggregated metrics, reporting            |
| Plugins    | Depends on plugin     | Optional storage per extension           |

**Principle:** Each microservice owns its database to ensure bounded context isolation.

---

## 9. Communication Patterns Summary

| Type           | Purpose                             | Example                                                     |
| -------------- | ----------------------------------- | ----------------------------------------------------------- |
| gRPC (sync)    | Real-time requests between services | Gateway → OnboardingService, OrdersService → CatalogService |
| Events (async) | Decoupled workflows                 | `StoreProvisioned`, `PaymentCompleted`, `OrderPaid`         |
| HTTP/GraphQL   | Client-facing                       | Dashboard requests, API for merchant apps                   |

---

## 10. Security Considerations
- JWT tokens for authentication
- Role-based access enforced at gateway
- TLS/SSL at NGINX edge
- No external exposure of microservices
- Sandbox accounts for trial merchants
- Secrets management for payments and plugins
- Audit logs for critical actions

---

## 11. Key Design Principles
1. DDD First: Bounded contexts with clear ownership
2. Event-Driven Workflows: Asynchronous orchestration
3. Microservices Isolation: Independent deployable and scalable units
4. Gateway Edge: Aggregation, auth, client-facing entry point
5. Zero-Touch Onboarding: Quick merchant setup
6. Plugin-Friendly Architecture: Extensions integrate via events
7. Observability & Metrics: Centralized logging, monitoring, tracing

---

## 12. Production Considerations
- Multi-tenancy support (tenant IDs or separate schemas)
- Horizontal scaling of gateway and core services
- Database backups per service
- Rate limiting at NGINX + Gateway
- Distributed tracing for complex workflows
- Disaster recovery planning
- Secure storage and rotation of merchant credentials

---

## 13. Summary
The StoreForge architecture delivers:
- Rapid merchant onboarding
- Secure, modular microservices with clear bounded contexts
- Direct merchant payouts via M-Pesa and Stripe
- Event-driven orchestration for provisioning, payments, and plugins
- Production-ready considerations like observability, security, and multi-tenancy
