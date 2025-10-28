# 🏗️ StoreForge

**StoreForge** is a self-service, multi-tenant ecommerce builder that lets anyone design, configure, and deploy a fully functional online store in **under 10 minutes** — no coding required.

> **"Forge your store in minutes."**

---

## 🌍 Overview

StoreForge enables small businesses and entrepreneurs to instantly create online stores using a drag-and-drop visual builder, integrated payment systems (Stripe, M-Pesa), and real-time hosting on custom domains.

It’s built for speed, scalability, and simplicity — combining a **Go backend**, **Next.js frontend**, and **multi-tenant architecture**.

---

## ✨ Features

- 🧱 **Visual Drag-and-Drop Builder**
- 💳 **Multi-payment support** (Stripe, M-Pesa, PayPal)
- 🛍️ **Product & order management**
- 🎨 **Customizable themes & templates**
- 🌐 **Instant subdomain & custom domain binding**
- 🔒 **Tenant isolation and secure data model**
- 💼 **Premium tiers and template marketplace**
- 🚀 **Deploy live in minutes**

---

## 🧠 Tech Stack & Rationale

| Layer | Tech | Why |
|-------|------|-----|
| Frontend | **Next.js with TS + Tailwind CSS** | Server-side rendering, React ecosystem, fast iteration |
| Builder Engine | **Craft.js / React DnD** | Drag-and-drop UI builder with JSON layout persistence |
| Backend | **Go (Gin + GORM)** | Concurrency, speed, and easy multi-tenant structure |
| Database | **PostgreSQL** | Reliability, relational data + JSONB support |
| Storage | **S3 / Cloudflare R2** | Cheap, scalable media storage |
| Proxy | **Caddy / Traefik** | Auto SSL and subdomain routing |
| Payments | **Stripe + M-Pesa** | Global + regional payment support |
| Caching | **Redis** | Tenant configuration and session caching |
| Deployment | **Docker + GitHub Actions** | Seamless CI/CD, reproducible environments |

---

## 🏗️ System Architecture

```
[User Browser]
   |
   v
[Next.js Frontend + Builder]  ---> calls ---> [Go API (Gin)]
                                              |
                                              v
                                      [PostgreSQL DB]
                                              |
                                              v
                                        [S3 / R2 Storage]
                                              |
                                              v
                                      [Payment APIs: Stripe/M-Pesa]
```

---

## 🧩 Folder Structure

```
storeforge/
 ├─ apps/
 │   ├─ frontend-next/       # Next.js storefront + builder + admin
 │   └─ api-go/              # Go backend (Gin)
 ├─ packages/
 │   ├─ ui-components/       # Shared React components
 │   └─ schemas/             # JSON schema definitions
 └─ infra/
     ├─ docker-compose.yml
     └─ k8s/
```

---

## ⚙️ Local Setup

### 1. Clone the repo
```bash
git clone https://github.com/yourusername/storeforge.git
cd storeforge
```

### 2. Start backend + frontend
```bash
docker-compose up
```

- Frontend: `http://localhost:4000`
- API: `http://localhost:8080`

### 3. Environment Variables
Create `.env` files for both apps:

**Frontend**
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

**Backend**
```env
DATABASE_URL=postgres://user:pass@db:5432/storeforge
STRIPE_SECRET_KEY=sk_test_xxx
MPESA_CONSUMER_KEY=xxx
MPESA_CONSUMER_SECRET=xxx
JWT_SECRET=supersecret
```

---

## 🧠 Core Concepts

### Multi-Tenancy
Each user (tenant) gets:
- A dedicated subdomain (`tenant.storeforge.io`)
- Isolated data using `tenant_id`
- Separate payment credentials

### Builder Engine
- Built using `Craft.js` or `React DnD`
- Saves page structure as JSON
- Renders published layout on storefront dynamically

### Security
- JWT authentication (per tenant)
- HTTPS enforced via Caddy
- Tenant-level access control
- Encrypted environment secrets

---

## 🪙 Premium Tiers

| Plan | Features |
|------|-----------|
| Free | 20 products, basic templates, subdomain only |
| Pro | Custom domain, advanced templates, no ads |
| Business | Analytics, multi-user, premium support |

---

## 🧩 API Overview

| Endpoint | Description |
|-----------|--------------|
| `/api/tenants` | Create, update tenant info |
| `/api/products` | CRUD products |
| `/api/orders` | Checkout flow |
| `/api/discounts` | Coupon management |
| `/api/layouts` | Save and publish layout schema |
| `/api/domains` | Custom domain management |

---

## 🧱 Roadmap

- [x] Tenant onboarding & subdomain setup  
- [x] Basic product CRUD  
- [ ] Drag-and-drop builder MVP  
- [ ] Stripe & M-Pesa integration  
- [ ] Publish + custom domain routing  
- [ ] Premium subscriptions  
- [ ] Template marketplace  
- [ ] Analytics dashboard  

---

## 🧪 Testing

- **Unit tests:** Go services and APIs  
- **Integration tests:** Builder + checkout flow  
- **E2E tests:** Frontend with Playwright or Cypress  

---

## 📊 KPIs

- 90% of users publish within 15 minutes  
- 99.9% uptime  
- <1% data leakage  
- 80% of live stores process at least 1 order  

---

## 🧭 License

MIT License © 2025 StoreForge

---

## 🧑‍💻 Contributors

- **Martin Wachira** – Founder & Lead Engineer  
  - Backend (Go), Infrastructure, System Architecture  
  - [LinkedIn](https://linkedin.com) | [Twitter](https://twitter.com)

Contributions, issues, and feature requests are welcome!

---

## 🌟 Support

If you like this project, consider giving it a ⭐ on GitHub and sharing it with your network!
