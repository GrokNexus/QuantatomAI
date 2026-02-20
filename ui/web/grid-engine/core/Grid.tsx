import { useGridQuery } from "../query/useGridQuery"
import { VirtualizedGrid } from "./VirtualizedGrid"

interface GridProps {
  query: any;
}

export function Grid({ query }: GridProps) {
  const { isLoading } = useGridQuery(query)

  if (isLoading) return <div>Loading grid…</div>

  return (
    <VirtualizedGrid
      rows={[]}
      columns={[]}
      cells={{}}
    />
  )
}
