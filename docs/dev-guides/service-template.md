# <ServiceName> Service

## Purpose

- **Domain:** e.g., Modeling / Planning / Actuals / Reconciliation / ALM / Connectors  
- **Responsibilities:**
  - Owns <domain> entities and APIs
  - Emits/consumes AODL events
  - Participates in DFN (Decentralized Flow Nexus)

## Tech stack

- Language: Go / Node.js / Rust (choose one per service)
- API: gRPC + REST (OpenAPI)
- Storage: Postgres (metadata) + AODL events
- Observability: OpenTelemetry (traces, metrics, logs)

## Layout

```text
<ServiceName>-service/
  src/
    main.(go|ts|rs)
    handlers/
    domain/
    adapters/
    config/
  tests/
  Dockerfile
  Makefile
  README.md
```
