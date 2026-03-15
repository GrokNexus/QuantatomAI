# Layer 4.2 Benchmark Evidence Template

## Purpose
This template standardizes how Phase 4 benchmark and replay evidence is recorded so results remain comparable across sessions and release candidates.

## Minimum Evidence Package
- run id
- git hash
- profile id
- environment summary
- governance controls enabled
- threshold targets
- observed latency and throughput metrics
- replay and correctness outcomes
- pass or fail decision

## Environment Summary
- execution date:
- operator:
- branch:
- git hash:
- environment type:
- database target:
- tenant count:
- duration seconds:

## Profile Summary
- profile id:
- workload class:
- command used:
- fixture dataset:

## Governance Controls Enabled
- phase 2 tenant controls:
- phase 3 metadata audit:
- workflow lock checks:
- ingest governance:

## Observed Metrics
- p50 latency ms:
- p95 latency ms:
- p99 latency ms:
- read throughput ops/sec:
- write throughput ops/sec:
- queue depth peak:
- tenant fairness ratio:
- audit amplification ratio:
- replay invalid rows:

## Threshold Evaluation
- threshold 1:
- threshold 2:
- threshold 3:
- overall result:

## Notes
- operational issues:
- anomalies observed:
- follow-up actions:

## Working Rule
Attach the evidence bundle created by [tools/load-testing/README.md](tools/load-testing/README.md) workflow and do not report benchmark claims without both the summary and raw command output.