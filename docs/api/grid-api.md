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
