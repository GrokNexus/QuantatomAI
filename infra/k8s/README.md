# QuantatomAI Kubernetes Guide

This guide covers local development and operations for the QuantatomAI platform using Kubernetes.

## Local Development (Kind)

To start a local cluster:

```bash
kind create cluster --name quantatom-local
```

### Applying Base Manifests

```bash
kubectl apply -k infra/k8s/base
```

### Core Components

- **AODL**: Append-Only Delta Log for eventing.
- **Hot Store**: Redis for low-latency lattice access.
- **Warm Store**: ClickHouse for analytical queries.
- **Domain Services**: Modeling and Planning.
- **UI**: Web dashboard.

## Multi-Cloud Overlays

Select the appropriate overlay for cloud deployment:

```bash
kubectl apply -k infra/k8s/overlays/aws
```
