# Grid API

## POST /grid/query

Executes a multidimensional grid query.

### Request

```json
{
  "dimensions": {
    "rows": ["Entity", "Product"],
    "columns": ["Time"],
    "pages": ["Scenario"],
    "filters": { "Region": ["NA"] }
  },
  "members": {
    "Entity": ["E100", "E200"],
    "Product": ["P100", "P200"],
    "Time": ["2025M01", "2025M02"],
    "Scenario": ["Working"]
  }
}
```

### Response

```json
{
  "rows": [{"id": "E100_P100", "name": "E100 - P100"}],
  "columns": [{"id": "2025M01", "name": "Jan 2025"}],
  "cells": [{"value": 1250.50}]
}
```

## POST /grid/writeback

Writes cell edits.

### Request

```json
{
  "cellEdits": [
    {
      "dims": { "Entity": "E100", "Product": "P200", "Time": "2025M01" },
      "measure": "Revenue",
      "scenario": "Working",
      "value": 12345.67
    }
  ]
}
```

### Response

```json
{
  "status": "ok"
}
```

## GET /api/v1/metadata/graph

Returns a tenant-scoped metadata lineage graph snapshot for a selected root member.

### Required Header

`X-Tenant-ID: <tenant-id>`

### Query Parameters

- `appId` (required)
- `dimension` (required)
- `rootMember` (required)

### Response

```json
{
  "tenantId": "tenant-ultra",
  "appId": "app-1",
  "dimension": "region",
  "rootMember": "Global",
  "nodes": [
    {"id": "global", "dimension": "region", "name": "Global", "path": "Global"},
    {"id": "global_northamerica", "dimension": "region", "name": "NorthAmerica", "path": "Global.NorthAmerica"}
  ],
  "edges": [
    {"fromId": "global", "toId": "global_northamerica", "type": "parent-child"}
  ],
  "ancestors": [],
  "descendants": ["NorthAmerica", "EMEA"]
}
```

## POST /api/v1/fluxion/forecast

Tenant-governed AI endpoint. Requires `X-Tenant-ID`. Returns `403` unless the tenant has explicitly opted in.

## POST /api/v1/fluxion/ask

Tenant-governed AI endpoint. Requires `X-Tenant-ID`. Returns `403` unless the tenant has explicitly opted in.
