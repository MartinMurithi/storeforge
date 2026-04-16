# Authentication & Tenant Bootstrap Flow

This document describes the **authoritative authentication, verification, and tenant bootstrap flow** for the StoreForge platform.

The design enforces **strong security boundaries**, **clear authorization context**, and **correct multi-tenant behavior**.

---

## Core Principles

1. **Identity is created before authorization**
2. **Tenants define authorization boundaries**
3. **All roles are tenant-scoped**
4. **Only tenant owners can onboard other users**
5. **JWTs are issued only when a tenant context exists**

---

## Phase 1: Registration (Identity Creation)

### Goal

Create a **user identity only**. No permissions, no tenant, no access token.

### Endpoint

```
POST /auth/register
```

### What Happens

* Validate input (email, phone, password, etc.)
* Hash password
* Create user record
* Set:

  * `is_verified = false`

### Explicitly NOT Done

* ❌ No tenant creation
* ❌ No role assignment
* ❌ No JWT issued

### Resulting State

```
User exists
User is unverified
User has no permissions
```

This phase answers **"Who are you?"**, not **"What can you do?"**

---

## Phase 2: Login + OTP (Authentication & Verification)

### Goal

Prove identity and verify ownership of contact details.

### Step 1: Login

```
POST /auth/login
```

#### What Happens

* Validate email & password
* If credentials are valid:

  * Generate OTP
  * Send OTP (SMS / Email)
  * Create login challenge

#### Response

```
{
  "login_challenge_id": "uuid"
}
```

No JWT is issued here.

---

### Step 2: Verify OTP

```
POST /auth/verify-otp
```

#### What Happens

* Validate OTP against challenge
* Mark user as:

  * `is_verified = true`
* Issue a **short-lived auth session** (internal use only)

### Resulting State

```
User is authenticated
User is verified
User still has no tenant
```

At this point:

* The system trusts the user’s identity
* The user still has **zero authorization**

---

## Phase 3: Tenant Bootstrap (First-Time Only)

### Critical Rule

> **Only verified users without an existing tenant may create a tenant.**

This is a **one-time bootstrap operation**.

### Check

```
Does user belong to any tenant?
```

* If **YES** → skip this phase
* If **NO** → tenant creation is mandatory

---

### Endpoint

```
POST /tenants
```

### Transactional Steps

All steps below happen **inside a single database transaction**:

1. Create tenant
2. Create `user_tenants` entry

   * `role = owner`
3. Set this tenant as user’s default tenant

### Resulting State

```
User exists
Tenant exists
User is owner of tenant
Authorization context is now defined
```

---

## Ownership Model (Important)

### Automatic Role Assignment

* **Only one role is ever assigned automatically:**

```
role = owner
```

This happens **only during tenant bootstrap**.

---

### Role Management Thereafter

* Owners may:

  * Invite users
  * Assign roles (`admin`, `staff`, etc.)
  * Remove users

* Non-owners:

  * ❌ Cannot invite users
  * ❌ Cannot assign roles
  * ❌ Cannot create tenants

This ensures **strict privilege control**.

---

## Phase 4: Issue Tenant-Scoped JWT (Final)

### Goal

Issue a JWT that is:

* Tenant-aware
* Role-aware
* Safe for API authorization

### When This Happens

Only after:

* User is verified
* Tenant exists
* Role is known

---

### JWT Payload Example

```json
{
  "sub": "user-id",
  "tid": "tenant-id",
  "role": "owner",
  "iss": "auth.storeforge",
  "aud": "storeforge-api",
  "exp": 1735689600
}
```

### Properties

* **Tenant-scoped** (`tid`)
* **Role-aware** (`role`)
* **Signed using RS256**
* Used for all authenticated API requests

---

## Authorization Model Summary

| Action        | Allowed          | Reason                       |
| ------------- | ---------------- | ---------------------------- |
| Register user | ✅                | Identity creation only       |
| Login         | ✅                | Authentication               |
| Verify OTP    | ✅                | Ownership verification       |
| Create tenant | ✅ (once)         | Bootstrap ownership          |
| Issue JWT     | ✅ (after tenant) | Authorization context exists |
| Add users     | ❌ (non-owner)    | Owner-only privilege         |
| Assign roles  | ❌ (non-owner)    | Prevent escalation           |

---

## Why This Flow:

1. Prevents premature authorization
2. Eliminates orphan tenants
3. Enforces least privilege
4. Supports clean multi-tenancy
5. Scales to multiple tenants per user

---

## Final Mental Model

```
Identity → Verification → Ownership → Authorization
```

Only after ownership exists does **power** (roles + JWT) get issued.

## TO-DO
1. FIX FOLDER STRUCTURE
2. IMPLEMENT EMAIL VERIFICATION
3. IMPLEMENT REFRESH TOKEN FLOW
4. LOGOUT FLOW
5. PASSWORD RESET FLOW
6. PBAC FLOW

```
user-management/
│
├─ cmd/                          # Entry point
│   └─ server/
│       └─ main.go               # Starts gRPC & HTTP server
│
├─ config/                       # Configuration files
│   └─ config.go
│
├─ bootstrap/                     # App initialization
│   └─ bootstrap.go
│
├─ proto/                         # gRPC protobuf definitions
│   └─ user.proto
│
├─ internal/
│   ├─ domain/
│   │   ├─ user.go                # Core User entity
│   │   ├─ role.go                # Roles & policies
│   │   └─ session.go             # Refresh token entity
│   │
│   ├─ repository/
│   │   ├─ user_repository.go     # DB operations for users
│   │   ├─ session_repository.go  # DB operations for refresh tokens
│   │   └─ role_repository.go     # Role storage / policies
│   │
│   ├─ service/
│   │   ├─ auth_service.go        # Registration, login, OTP verification
│   │   ├─ session_service.go     # JWT + refresh token logic
│   │   ├─ password_service.go    # Reset, change password
│   │   └─ pbac_service.go        # Policy evaluation & enforcement
│   │
│   ├─ handlers/                  # HTTP / gRPC endpoints
│   │   ├─ auth_handler.go
│   │   ├─ session_handler.go
│   │   └─ password_handler.go
│   │
│   ├─ dto/                       # Input/Output structures
│   │   ├─ auth_dto.go
│   │   ├─ session_dto.go
│   │   └─ password_dto.go
│   │
│   ├─ mappers/                   # Map DTOs to domain entities and vice versa
│   │   └─ user_mapper.go
│   │
│   ├─ errors/                    # Custom error types
│   │   ├─ app_errors.go
│   │   └─ db_errors.go
│   │
│   └─ utils/                     # Helper functions
│       ├─ jwt_utils.go
│       ├─ otp_utils.go
│       └─ hash_utils.go
│
├─ database/
│   ├─ migrations/
│   ├─ config/
│   └─ seeding/
│
└─ go.mod
```