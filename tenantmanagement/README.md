````markdown
# StoreForge Schema-Driven Runtime Engine

## 🧠 Overview

StoreForge is a multi-tenant ecommerce platform designed around a **schema-driven runtime architecture**. Instead of hardcoding storefronts, admin behavior, and UI structure, StoreForge uses a **strict, versioned configuration system** that defines how each tenant’s store looks and behaves.

At its core, StoreForge is not just an ecommerce backend—it is a:

> **Schema-Driven Commerce Runtime Engine (SDCRE)**

This means the system interprets structured configuration (JSON schemas) at runtime to dynamically generate:
- Storefront UI (customer-facing experience)
- Admin dashboard behavior (merchant-facing system)
- Business workflows (orders, payments, products)
- Theme and design systems per tenant

---

## ⚙️ Core Idea

Traditional ecommerce systems are **code-driven**:
- UI is hardcoded in frontend
- Business logic is embedded in backend
- Customization requires engineering changes

StoreForge is **configuration-driven**:
- UI is defined by schema
- Behavior is defined by structured actions
- Backend enforces rules, frontend renders dynamically
- Each tenant becomes a “runtime instance” of the system

---

## 🧩 System Components

### 1. Theme System (Design Layer)

Defines the visual identity of a store using design tokens.

Includes:
- Colors (primary, secondary, background, text)
- Typography scale
- Spacing system
- Component styling overrides

This is NOT raw CSS. It is a **design token abstraction layer**.

---

### 2. Storefront Schema (UI Composition Layer)

Defines the customer-facing website structure.

Key concepts:
- Pages (home, product, catalog, cart)
- Sections (hero, product grid, banners, etc.)
- Layout rules (grid, list, responsive behavior)
- Product display logic (filters, sorting, limits)

Each page is a **composable tree of components defined in JSON**.

Frontend maps schema → React components via a registry system.

---

### 3. Admin Schema (Operations Layer)

Defines the merchant dashboard behavior per tenant.

Includes:
- Modules (products, orders, users, payments)
- Feature toggles per tenant
- Role-based access control (RBAC)
- Order workflow state machines (e.g. pending → paid → shipped)

This ensures that even admin behavior is **configurable and tenant-specific**.

---

### 4. Action System (Interaction Layer)

Defines how UI elements behave.

Instead of embedding logic in UI code, elements define actions:

Example:
```json
{
  "type": "button",
  "props": {
    "label": "Add to Cart",
    "action": {
      "type": "add_to_cart",
      "payload": {
        "product_id": "$product.id"
      }
    }
  }
}
````

The frontend interprets these actions via a controlled execution layer.

This ensures:

* No arbitrary code execution from JSON
* Safe, predictable interactivity
* Centralized business logic handling

---

### 5. Versioning System

All schemas are versioned using semantic versioning:

* `MAJOR` → breaking schema changes
* `MINOR` → backward-compatible additions
* `PATCH` → fixes or corrections

Each tenant config includes:

```json
{
  "schema_version": "1.2.0"
}
```

Backend handles:

* validation
* migration
* backward compatibility transformation

---

## 🧱 Architecture Model

StoreForge follows a strict layered model:

```
[ Theme Layer ]        → visual design system
[ Storefront Layer ]   → customer UI structure
[ Admin Layer ]        → merchant operations system
[ Action Layer ]       → interaction behavior model
[ Engine Layer ]       → validation + interpretation runtime
```

---

## ⚙️ Runtime Flow

### 1. Backend (Go)

* Stores tenant configuration
* Validates schema strictly
* Applies version migrations
* Serves config via API/gRPC

### 2. Frontend (Next.js)

* Fetches tenant schema
* Uses component registry to render UI
* Executes actions via controlled handlers
* Applies theme tokens dynamically

### 3. Execution Model

1. JSON schema is loaded
2. Schema is validated
3. Components are resolved via registry
4. UI is rendered dynamically
5. User actions are mapped to predefined handlers

---

## 🔐 Design Constraints

To ensure system safety and scalability:

* ❌ No free-form JSON execution
* ❌ No arbitrary JavaScript in config
* ❌ No dynamic unknown components
* ✅ Strict schema validation in backend
* ✅ Whitelisted UI component registry
* ✅ Controlled action execution system
* ✅ Tenant isolation enforced at all layers

---

## 🚀 Why This Architecture Matters

This system enables:

### 1. True Multi-Tenancy

Each tenant behaves like an isolated “application instance”.

### 2. Dynamic UI Without Redeployments

Stores can change layout, structure, and content without frontend changes.

### 3. Platform Scalability

New features can be added via schema extensions instead of code rewrites.

### 4. Enterprise-Level Customization

Tenants can have:

* unique storefront layouts
* custom workflows
* personalized admin systems

---

## 🧠 Mental Model

Think of StoreForge as:

> A runtime interpreter for ecommerce applications

Not a traditional monolithic backend.

Instead:

* JSON = program definition
* Engine = interpreter
* Frontend = rendering runtime
* Backend = validation + execution authority

---

## 📌 Final Summary

StoreForge is a **Schema-Driven Commerce Runtime Engine** that transforms ecommerce from a hardcoded system into a **declarative, configurable platform**.

It enables:

* dynamic storefront generation
* tenant-specific customization
* strict backend enforcement
* safe, structured UI composition

At scale, this allows StoreForge to function as a **multi-tenant ecommerce operating system**, not just a storefront builder.
