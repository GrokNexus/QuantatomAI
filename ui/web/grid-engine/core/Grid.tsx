"use client";

import { useGridQuery, GridQueryPayload } from "../query/useGridQuery"
import { VirtualizedGrid } from "./VirtualizedGrid"

interface GridProps {
  query?: GridQueryPayload;
}

export function Grid({ query }: GridProps) {
  const { data, isLoading, error } = useGridQuery(query)

  if (isLoading) return <div>Loading grid…</div>

  if (error) return <div>Grid error: {error}</div>

  if (!data) return <div>No grid data</div>

  return (
    <VirtualizedGrid
      rows={data.rows}
      columns={data.columns}
      cells={data.cells}
    />
  )
}
