# Grid Operations Spec Sheet

## Parity & Beyond
| Operation | Current SLO (p95) | Target SLO (p95) | Notes |
|-----------|-------------------|------------------|-------|
| Grid Open (Hot) | 700ms | 400ms | Cache warming in progress. |
| Grid Open (Cold) | 1.5s | 800ms | Background fetching enabled. |
| Single Cell Edit | 600ms | 300ms | Optimistic UI updates. |
| Medium Spread | 1.2s | 1.0s | Parallel calculation. |
| Heavy Allocation | 10s | 5s | Job sizing optimization. |
| Actuals Load | 10 min | 2 min | Ingestion pipeline tuning. |
