
# StoreForge

**StoreForge** is a microservices-based Platform-as-a-Service (PaaS) that enables merchants to launch fully functional e-commerce stores in under **5 minutes**. Built with **Go, gRPC, envoy proxy event-driven architecture, and modular microservices**, it supports trial-first subscriptions, direct merchant payouts, and a plugin-friendly ecosystem for extensibility.

This repository serves as the reference for the StoreForge system architecture, design patterns, and technical implementation.

---

## Table of Contents

- [StoreForge](#storeforge)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Architecture Overview](#architecture-overview)
    - [Architectural Principles](#architectural-principles)
    - [Communication Patterns](#communication-patterns)
  - [Core Services](#core-services)
    - [Onboarding Service](#onboarding-service)
    - [Auth Service](#auth-service)
  - [Supporting Services](#supporting-services)
    - [Catalog Service](#catalog-service)
    - [Payments Service](#payments-service)
    - [Storefront / Provisioning Service](#storefront--provisioning-service)
    - [Notifications Service](#notifications-service)
    - [Analytics / Observability Service](#analytics--observability-service)
    - [Subscription / Feature Flags Service](#subscription--feature-flags-service)
    - [Extensions / Plugin Registry Service](#extensions--plugin-registry-service)
  - [Operational / Supporting Entities](#operational--supporting-entities)
    - [Orders Service](#orders-service)
    - [Customers Service](#customers-service)
  - [Databases \& Storage Strategy](#databases--storage-strategy)
  - [Security Considerations](#security-considerations)
  - [Key Design Principles](#key-design-principles)
  - [Production Considerations](#production-considerations)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Running Locally](#running-locally)

---

## Introduction

StoreForge is designed for **rapid merchant onboarding**, **scalable microservices**, and **direct merchant payouts** via platforms like **M-Pesa** and **Stripe**.  

**Key Goals:**
- Enable merchants to start selling in minutes  
- Ensure secure handling of payments and sensitive data  
- Provide modular, extensible architecture for plugins and additional services  
- Maintain production-grade observability, metrics, and audit trails  

---

## Architecture Overview

### Architectural Principles

- **Domain-Driven Design (DDD):** Each bounded context owns its data and business logic  
- **Microservices Isolation:** Services are independently deployable and scalable  
- **Event-Driven Workflows:** Asynchronous events orchestrate business processes  
- **Gateway Edge:** Single entry point for clients, enforcing authentication and role-based access  
- **Zero-Touch Onboarding:** Minimal merchant involvement to launch stores  
- **Plugin/Extension Ecosystem:** Optional features integrate via events  

### Communication Patterns

| Type           | Purpose                             | Example                                                     |
| -------------- | ----------------------------------- | ----------------------------------------------------------- |
| gRPC (sync)    | Real-time requests between services | Gateway → OnboardingService, OrdersService → CatalogService |
| Events (async) | Decoupled workflows                 | `StoreProvisioned`, `PaymentCompleted`, `OrderPaid`         |
| HTTP | Client-facing API                   | Dashboard requests, API for merchant apps                   |

---

## Core Services

### Onboarding Service
- Manages **merchant lifecycle** and **store provisioning**  
- Handles URL assignment, theme setup, default product catalog creation  
- Publishes events: `MerchantVerified`, `StoreProvisioned`  

### Auth Service
- Handles **registration, login, OTP verification, password management**  
- Issues **JWTs** and manages refresh tokens  
- Implements **PBAC** for role-based authorization  

---

## Supporting Services

### Catalog Service
- Manages **products, categories, and inventory**  
- Publishes inventory change events  

### Payments Service
- Supports **M-Pesa**, and other plugins  
- Differentiates between **sandbox (trial)** and **live accounts**  
- Direct payouts to merchants, platform does not handle funds  
- Publishes events: `PaymentCompleted`, `PaymentFailed`, `PaymentRefunded`  

### Storefront / Provisioning Service
- Applies **themes, hosting configuration**, and store readiness  

### Notifications Service
- Sends **emails, SMS, OTPs**, and webhooks  
- Reacts to key events like `PaymentCompleted` or `StoreProvisioned`  

### Analytics / Observability Service
- Aggregates **metrics, logs, and performance data**  
- Provides dashboard information for merchants  

### Subscription / Feature Flags Service
- Manages **trial periods, plan upgrades, feature access**  

### Extensions / Plugin Registry Service
- Optional integrations like **Fleetbase delivery, advanced analytics**  
- Plugins integrate **asynchronously via events**

---

## Operational / Supporting Entities

### Orders Service
- Transaction hub for merchant sales  
- Publishes: `OrderCreated`, `OrderPaid`, `OrderShipped`, `OrderCompleted`  

### Customers Service
- Manages buyer profiles, purchase history, and optional segmentation  

---

## Databases & Storage Strategy

| Service    | Database Type         | Notes                                     |
| ---------- | --------------------- | ----------------------------------------- |
| Auth       | PostgreSQL            | User credentials, roles, JWTs             |
| Onboarding | PostgreSQL            | Merchant & store lifecycle                |
| Catalog    | PostgreSQL/MongoDB    | Products, categories, inventory           |
| Payments   | PostgreSQL            | Transaction ledger, payment status        |
| Orders     | PostgreSQL            | Order lifecycle & transactional integrity |
| Customers  | PostgreSQL/MongoDB    | Buyer data, purchase history              |
| Analytics  | Clickhouse/PostgreSQL | Aggregated metrics, reporting             |
| Plugins    | Depends on plugin     | Optional storage per extension            |

**Principle:** Each service owns its database to ensure **bounded context isolation**.

---

## Security Considerations

- JWT tokens for authentication  
- Policy-based access enforced at **gateway**  
- TLS/SSL termination at **NGINX edge**  
- Sandbox accounts for trial merchants  
- Secrets management for payments and plugins  
- Audit logs for critical actions  

---

## Key Design Principles

1. DDD First: Bounded contexts with clear ownership  
2. Event-Driven Workflows: Asynchronous orchestration  
3. Microservices Isolation: Independent deployable units  
4. Gateway Edge(Envoy): Aggregation, auth, client-facing entry point  
5. Zero-Touch Onboarding: Quick merchant setup  
6. Plugin-Friendly Architecture: Extensions via events  
7. Observability & Metrics: Centralized logging and monitoring  

---

## Production Considerations

- Multi-tenancy support (tenant IDs or separate schemas)  
- Horizontal scaling of gateway and core services  
- Database backups per service  
- Rate limiting at NGINX + Gateway  
- Distributed tracing for complex workflows  
- Disaster recovery planning  
- Secure storage and rotation of merchant credentials  

---

## Getting Started

### Prerequisites
- Go 1.20+  
- PostgreSQL && MongoDB  
- Message broker: NATS  
- Docker & Docker Compose

### Running Locally
1. Clone the repository  
   ```bash
   git clone https://github.com/MartinMurithi/storeforge.git
   cd storeforge
Start databases and message broker

bash
Copy code
docker-compose up -d
Run gRPC servers for core services (e.g., Auth, Onboarding, Catalog)

Start API Gateway

Access dashboard at http://localhost:8080