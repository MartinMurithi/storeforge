1️⃣ Envoy — Infrastructure Layer (Deep Explanation)

Envoy operates at Layer 7 networking but remains infrastructure-focused, not product-focused.

Think:

Envoy manages traffic behavior, not business meaning.

It should answer questions like:

Where should this request go?

Is this client allowed?

Is the backend healthy?

Should we retry?

How do we secure the connection?

Not:

What does a product look like?

How should pagination be formatted?

What fields should the frontend see?

Envoy Responsibilities in Detail
1. TLS Termination

Envoy handles HTTPS encryption.

Flow:

Client HTTPS → Envoy decrypts → internal HTTP/gRPC

Benefits:

Services don’t manage certificates

Centralized security

Easier rotation (LetsEncrypt, ACM, etc.)

2. Routing

Envoy decides where requests go.

Example:

/api/users → gateway-service
/grpc.products → product-service

Routing rules are declarative config.

3. Load Balancing

If you have multiple instances:

gateway-1
gateway-2
gateway-3

Envoy distributes traffic:

Round robin

Least request

Weighted

Health-aware

Your services stay simple.

4. Authentication (JWT Verification)

Envoy can validate tokens before requests reach your code.

Example:

Authorization: Bearer <JWT>

Envoy checks:

Signature

Expiration

Issuer

Audience

If invalid → request blocked.

This reduces load on services.

Important: Envoy verifies identity, but authorization logic still belongs in gateway.

5. Rate Limiting

Protect your system from abuse.

Examples:

100 requests/min per user

1000 requests/min per IP

Tier-based limits

Envoy integrates with rate-limit services.

6. Observability / Tracing

Envoy automatically emits:

Metrics (Prometheus)

Distributed tracing (Jaeger, Zipkin)

Access logs

You get system visibility without touching app code.

7. Retries / Circuit Breaking

Envoy handles resilience patterns:

Retry:

Service failed → retry another instance

Circuit breaker:

Service unhealthy → stop sending traffic

This prevents cascading failures.

Why Envoy Should NOT Shape JSON

Because:

Config becomes complex

Logic becomes untestable

Business rules leak into infra

Versioning becomes impossible

Engineers lose control

Envoy is not designed to be a product layer.