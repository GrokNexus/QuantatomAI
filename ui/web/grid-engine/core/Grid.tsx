import { useGridQuery } from "../query/useGridQuery"
import { VirtualizedGrid } from "./VirtualizedGrid"

export function Grid({ query }) {
  const { data, isLoading } = useGridQuery(query)

  if (isLoading) return <div>Loading grid…</div>

  return (
    <VirtualizedGrid
      rows={data.rows}
      columns={data.columns}
      cells={data.cells}
    />
  )
}
