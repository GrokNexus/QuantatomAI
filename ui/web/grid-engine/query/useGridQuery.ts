import { useEffect, useMemo, useState } from 'react';

interface MemberInfo {
  id: number;
  code: string;
  name: string;
}

interface CellValue {
  rowIndex: number;
  colIndex: number;
  value: number;
}

interface GridQueryResponse {
  rows: MemberInfo[][];
  columns: MemberInfo[][];
  cells: CellValue[];
}

export interface GridData {
  rows: string[];
  columns: string[];
  cells: Record<string, number>;
}

export interface GridQueryPayload {
  dimensions?: {
    rows?: string[];
    columns?: string[];
    pages?: string[];
    filters?: Record<string, string[]>;
  };
  members?: Record<string, string[]>;
  window?: {
    rowStart?: number;
    rowEnd?: number;
    colStart?: number;
    colEnd?: number;
  };
  defaults?: Record<number, number>;
  stream?: boolean;
  branchId?: string;
}

const EMPTY_QUERY: GridQueryPayload = {};

const DEFAULT_QUERY: GridQueryPayload = {
  dimensions: {
    rows: ['Entity'],
    columns: ['Time'],
    pages: [],
    filters: {},
  },
  members: {
    Entity: ['North America', 'EMEA', 'APAC'],
    Time: ['2025', '2026', '2027'],
    Measure: ['Revenue'],
    Scenario: ['Actual'],
  },
  window: { rowStart: 0, rowEnd: 0, colStart: 0, colEnd: 0 },
  defaults: {},
  stream: false,
};

export function useGridQuery(query?: GridQueryPayload) {
  const [data, setData] = useState<GridData | null>(null);
  const [rowCount, setRowCount] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const inputQuery = query ?? EMPTY_QUERY;

  const payload = useMemo(() => ({
    ...DEFAULT_QUERY,
    ...inputQuery,
    dimensions: { ...DEFAULT_QUERY.dimensions, ...inputQuery.dimensions },
    members: { ...DEFAULT_QUERY.members, ...inputQuery.members },
    window: { ...DEFAULT_QUERY.window, ...inputQuery.window },
    defaults: { ...DEFAULT_QUERY.defaults, ...inputQuery.defaults },
  }), [inputQuery]);

  useEffect(() => {
    let isMounted = true;
    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const endpoint = process.env.NEXT_PUBLIC_GRID_PROXY || '/api/grid';
        const res = await fetch(endpoint, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(payload),
        });

        if (!res.ok) {
          throw new Error(`Grid service responded with ${res.status}`);
        }

        const json = (await res.json()) as GridQueryResponse;

        if (!isMounted) return;

        const columns = (json.columns || []).map((combo, idx) => {
          const label = combo.map((m) => m.name || m.code).filter(Boolean).join(' / ');
          return label || `Col ${idx + 1}`;
        });

        const rows = (json.rows || []).map((combo, idx) => {
          const label = combo.map((m) => m.name || m.code).filter(Boolean).join(' / ');
          return label || `Row ${idx + 1}`;
        });

        const cells: Record<string, number> = {};
        (json.cells || []).forEach((cell) => {
          const key = `${cell.rowIndex}-${cell.colIndex}`;
          cells[key] = cell.value;
        });

        setData({ rows, columns, cells });
        setRowCount(rows.length);
      } catch (err: any) {
        console.error("Grid query failed", err);
        if (isMounted) {
          setError(err.message);

          // Fallback demo data to keep UI usable if backend is down
          const rows = ['North America', 'EMEA', 'APAC', 'LATAM'];
          const columns = ['Total Revenue', 'Hardware', 'Software', 'Services'];
          const fallbackCells: Record<string, number> = {};
          rows.forEach((_, r) => {
            columns.forEach((_, c) => {
              fallbackCells[`${r}-${c}`] = Math.round(Math.random() * 5_000_000) / 100;
            });
          });
          setData({ rows, columns, cells: fallbackCells });
          setRowCount(rows.length);
        }
      } finally {
        if (isMounted) setIsLoading(false);
      }
    };

    fetchData();

    return () => { isMounted = false; };
  }, [payload]);

  return { data, rowCount, isLoading, error };
}
