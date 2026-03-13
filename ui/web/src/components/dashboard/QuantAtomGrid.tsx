"use client";

// QuantAtomGrid.tsx ── unified grid engine

import React, { useState, useEffect, useRef, useMemo, useCallback } from 'react';
import {
  TrendingUp,
  TrendingDown,
  ArrowUpRight,
  Filter,
  Download,
  MoreHorizontal,
  Plus,
  Eye,
  EyeOff,
  Maximize2,
  Minimize2,
  X,
  RotateCcw,
  MessageCircle,
  Users,
  Copy,
  FileDown,
  AlertTriangle,
  Undo2,
  Redo2,
  BarChart as BarIcon,
  LineChart as LineIcon,
  GripVertical,
  GripHorizontal
} from 'lucide-react';

import { evaluate as mathEvaluate } from 'mathjs';
import * as XLSX from 'xlsx';
import {
  BarChart,
  Bar,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts';
import { DndProvider, useDrag, useDrop } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';

// ────────────────────────────────────────────────
// Types
// ────────────────────────────────────────────────

interface Dimension {
  name: string;
  members: string[];
}

interface RowData {
  id: string;
  [key: string]: any;
  children?: RowData[];
}

interface NamedRanges {
  [name: string]: string;
}

interface OnlineUser {
  id: string;
  name: string;
  avatar: string;
  color: string;
}

interface CellActivity {
  [key: string]: string; // `${rowId}-${colKey}` → userId
}

interface Comment {
  user: string;
  text: string;
  time: string;
}

interface CellComments {
  [key: string]: Comment[];
}

interface ValidationRule {
  type?: 'number' | 'text' | 'date' | 'custom';
  required?: boolean;
  min?: number;
  max?: number;
  minLength?: number;
  maxLength?: number;
  pattern?: RegExp;
  customMessage?: string;
}

interface EditHistoryEntry {
  timestamp: number;
  rowId: string;
  colKey: string;
  oldValue: any;
  newValue: any;
  hierarchySnapshot: RowData[];
}

type FormulaResult = string | number | boolean | any[] | null;

const ItemTypes = {
  COLUMN: 'column',
  ROW: 'row'
};

interface MatrixRow {
  month: string;
  [year: string]: string | number;
}

interface VectorProjectionRow {
  month: string;
  year: string;
  periodKey: string;
  yield: number;
  area: number;
  weight: number;
  price: number;
  revenue: number;
}

const MATRIX_MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'];
const MATRIX_YEARS = ['2025', '2026', '2027'];
const PRICE_BY_PERIOD: Record<string, number> = {
  '2025-Jan': 12.0,
  '2025-Feb': 12.2,
  '2025-Mar': 12.5,
  '2025-Apr': 12.7,
  '2025-May': 12.9,
  '2025-Jun': 13.1,
  '2026-Jan': 13.8,
  '2026-Feb': 14.0,
  '2026-Mar': 14.3,
  '2026-Apr': 14.5,
  '2026-May': 14.8,
  '2026-Jun': 15.0,
  '2027-Jan': 15.5,
  '2027-Feb': 15.8,
  '2027-Mar': 16.1,
  '2027-Apr': 16.3,
  '2027-May': 16.6,
  '2027-Jun': 16.9
};

// ────────────────────────────────────────────────
// computePivot
// ────────────────────────────────────────────────

function computePivot(
  data: RowData[],
  rowField: string,
  colField: string,
  valueField: string,
  agg: 'sum' | 'avg' | 'count' = 'sum'
) {
  const rowKeys = Array.from(new Set(data.map((r) => r[rowField])));
  const colKeys = Array.from(new Set(data.map((r) => r[colField])));
  const table: Record<string, Record<string, any[]>> = {};
  rowKeys.forEach((rk) => {
    table[rk] = {};
    colKeys.forEach((ck) => {
      table[rk][ck] = [];
    });
  });
  data.forEach((row) => {
    if (row[rowField] && row[colField] && row[valueField] !== undefined) {
      table[row[rowField]][row[colField]].push(row[valueField]);
    }
  });

  const aggFn = (arr: any[]) => {
    if (agg === 'sum') return arr.reduce((a: number, b: number) => a + Number(b), 0);
    if (agg === 'avg') return arr.length ? arr.reduce((a: number, b: number) => a + Number(b), 0) / arr.length : 0;
    if (agg === 'count') return arr.length;
    return '';
  };

  const result = rowKeys.map((rk) => colKeys.map((ck) => aggFn(table[rk][ck])));
  return { rowKeys, colKeys, result };
}

// ────────────────────────────────────────────────
// Formula Engine
// ────────────────────────────────────────────────

const parseArgs = (fnCall: string): string[] => {
  const match = fnCall.match(/\((.*)\)/);
  if (!match) return [];
  return match[1].split(',').map((a) => a.trim()).filter(Boolean);
};

const evaluateFormula = (
  formula: string,
  row: RowData,
  gridData: RowData[] = [],
  namedRanges: NamedRanges = {}
): FormulaResult => {
  if (!formula.startsWith('=')) return formula;

  let expr = formula.slice(1).trim();

  Object.entries(namedRanges).forEach(([name, value]) => {
    const regex = new RegExp(`\\b${name}\\b`, 'gi');
    expr = expr.replace(regex, value);
  });

  Object.entries(row).forEach(([key, value]) => {
    const regex = new RegExp(`\\b${key}\\b`, 'gi');
    expr = expr.replace(regex, String(value ?? 0));
  });

  // Example statistical functions (add more as needed)
  if (/^SUM\(/i.test(expr)) {
    const args = parseArgs(expr);
    try {
      const arr: number[] = JSON.parse(args[0]);
      return arr.reduce((a: number, b: number) => a + b, 0);
    } catch {
      return 'ERR';
    }
  }

  // Add your other functions here (VAR, STDEV, IRR, etc.) from previous versions

  try {
    return mathEvaluate(expr, { ...row });
  } catch (err) {
    return `ERR: ${(err as Error).message}`;
  }
};

// ────────────────────────────────────────────────
// Validation Helper
// ────────────────────────────────────────────────

const validateCell = (value: any, rule?: ValidationRule): string => {
  if (!rule) return '';

  const strVal = String(value ?? '');

  if (rule.required && (value === undefined || value === null || strVal === '')) {
    return rule.customMessage || 'Required';
  }

  if (rule.type === 'number') {
    const num = Number(value);
    if (isNaN(num)) return 'Must be a number';
    if (rule.min !== undefined && num < rule.min) return `Minimum ${rule.min}`;
    if (rule.max !== undefined && num > rule.max) return `Maximum ${rule.max}`;
  }

  if (rule.type === 'text') {
    if (rule.minLength !== undefined && strVal.length < rule.minLength) return `Minimum length ${rule.minLength}`;
    if (rule.maxLength !== undefined && strVal.length > rule.maxLength) return `Maximum length ${rule.maxLength}`;
    if (rule.pattern && !rule.pattern.test(strVal)) return 'Invalid format';
  }

  if (rule.type === 'date') {
    const date = new Date(value);
    if (isNaN(date.getTime())) return 'Invalid date';
  }

  return '';
};

// ────────────────────────────────────────────────
// Main Component
// ────────────────────────────────────────────────

interface QuantAtomGridProps {
  appId?: string;
  cubeEndpoint?: string;
  dimensions?: Dimension[];
  columnOrder?: string[];
}

export function QuantAtomGrid({
  appId = 'QuantAtom Planning Grid',
  cubeEndpoint = '/api/olap/sales',
  dimensions: initialDimensions = [
    { name: 'Time', members: ['2025', '2026', '2027', 'All'] },
    { name: 'Entity', members: ['North America', 'EMEA', 'APAC', 'All'] }
  ],
  columnOrder: initialColumnOrder = ['entity', 'time', 'budget', 'actual', 'variance']
}: QuantAtomGridProps) {
  const [dimensions] = useState<Dimension[]>(initialDimensions);
  const [selectedMembers, setSelectedMembers] = useState<Record<string, string>>({});
  const [hierarchyData, setHierarchyData] = useState<RowData[]>([]);
  const [loading, setLoading] = useState(false);

  const [theme, setTheme] = useState('default');
  const themeOptions = ['default', 'dark', 'glass', 'emerald', 'rose', 'amber'];

  const [expandedTreeRows, setExpandedTreeRows] = useState<string[]>([]);
  const [columnOrder, setColumnOrder] = useState<string[]>(initialColumnOrder);
  const [hiddenColumns, setHiddenColumns] = useState<string[]>([]); // Added stub to fix reference error

  // Editing
  const [editingCell, setEditingCell] = useState<{ rowId: string; colKey: string; value: string } | null>(null);

  // Pivoting
  const [pivotMode, setPivotMode] = useState(false);
  const [pivotResult, setPivotResult] = useState<{ rowKeys: string[]; colKeys: string[]; result: any[][] } | null>(null);

  // Collaboration
  const [cellActivity, setCellActivity] = useState<CellActivity>({});
  const [onlineUsers, setOnlineUsers] = useState<OnlineUser[]>([]);
  const [cellComments, setCellComments] = useState<CellComments>({});
  const [newComment, setNewComment] = useState('');
  const [commentCell, setCommentCell] = useState<{ rowId: string; colKey: string } | null>(null);
  const [collabPanelOpen, setCollabPanelOpen] = useState(false);
  const [commentsPanelOpen, setCommentsPanelOpen] = useState(false);

  // Validation
  const [validationRules, setValidationRules] = useState<Record<string, ValidationRule>>({});
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});
  const [showValidationPanel, setShowValidationPanel] = useState(false);

  // Undo/Redo
  const [history, setHistory] = useState<EditHistoryEntry[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const MAX_HISTORY = 50;

  // Charting
  const [showChart, setShowChart] = useState(false);
  const [chartType, setChartType] = useState<'bar' | 'line'>('bar');

  // Column resize
  const [columnWidths, setColumnWidths] = useState<Record<string, number>>({});

  // Integrated Matrix/Vector capability (inside the same grid engine)
  const [gridStyle, setGridStyle] = useState<'matrix' | 'vector'>('vector');
  const [yieldMatrix, setYieldMatrix] = useState<MatrixRow[]>(
    MATRIX_MONTHS.map((month) => ({ month, '2025': 0, '2026': 0, '2027': 0 }))
  );
  const [areaMatrix, setAreaMatrix] = useState<MatrixRow[]>(
    MATRIX_MONTHS.map((month) => ({ month, '2025': 0, '2026': 0, '2027': 0 }))
  );

  const weightMatrix = useMemo(() => {
    return MATRIX_MONTHS.map((month, monthIndex) => {
      const row: MatrixRow = { month };
      MATRIX_YEARS.forEach((year) => {
        const yieldValue = Number(yieldMatrix[monthIndex][year]) || 0;
        const areaValue = Number(areaMatrix[monthIndex][year]) || 0;
        row[year] = yieldValue * areaValue;
      });
      return row;
    });
  }, [yieldMatrix, areaMatrix]);

  const vectorProjectionRows = useMemo<VectorProjectionRow[]>(() => {
    const rows: VectorProjectionRow[] = [];
    MATRIX_MONTHS.forEach((month, monthIndex) => {
      MATRIX_YEARS.forEach((year) => {
        const periodKey = `${year}-${month}`;
        const yieldValue = Number(yieldMatrix[monthIndex][year]) || 0;
        const areaValue = Number(areaMatrix[monthIndex][year]) || 0;
        const weightValue = yieldValue * areaValue;
        const priceValue = PRICE_BY_PERIOD[periodKey] ?? 0;
        const revenueValue = weightValue * priceValue;
        rows.push({
          month,
          year,
          periodKey,
          yield: yieldValue,
          area: areaValue,
          weight: weightValue,
          price: priceValue,
          revenue: revenueValue
        });
      });
    });
    return rows;
  }, [yieldMatrix, areaMatrix]);

  const revenueMatrix = useMemo(() => {
    const revenueByPeriod = new Map<string, number>(
      vectorProjectionRows.map((row) => [row.periodKey, row.revenue])
    );
    return MATRIX_MONTHS.map((month, monthIndex) => {
      const row: MatrixRow = { month };
      MATRIX_YEARS.forEach((year) => {
        const periodKey = `${year}-${month}`;
        row[year] = revenueByPeriod.get(periodKey) ?? 0;
      });
      return row;
    });
  }, [vectorProjectionRows]);

  const handleMatrixInput = (
    setter: React.Dispatch<React.SetStateAction<MatrixRow[]>>,
    monthIndex: number,
    year: string,
    value: string
  ) => {
    setter((prev) => {
      const next = [...prev];
      next[monthIndex] = { ...next[monthIndex], [year]: Number(value) || 0 };
      return next;
    });
  };

  const gridRef = useRef<HTMLDivElement>(null);

  // Fetch data
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const params = new URLSearchParams(selectedMembers);
        const url = `${cubeEndpoint}?${params.toString()}`;
        const res = await fetch(url);
        if (!res.ok) throw new Error('Fetch failed');
        const json = await res.json();
        setHierarchyData(json.data || json || []);
      } catch (err) {
        console.error('OLAP fetch error:', err);
        setHierarchyData([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [selectedMembers, cubeEndpoint]);

  // ... (rest of the hooks and logic from previous steps - online users, validation, undo/redo, etc.)

  // Export CSV
  const exportToCSV = () => {
    const flatData = flattenHierarchy(hierarchyData);
    const headers = columnOrder.join(',');
    const rows = flatData.map(row => columnOrder.map(col => row[col] ?? '').join(','));
    const csv = [headers, ...rows].join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${appId.replace(/\s+/g, '_')}_export.csv`;
    a.click();
    URL.revokeObjectURL(url);
  };

  // Export Excel
  const exportToExcel = () => {
    const flatData = flattenHierarchy(hierarchyData);
    const ws = XLSX.utils.json_to_sheet(flatData, { header: columnOrder });
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Data');
    XLSX.writeFile(wb, `${appId.replace(/\s+/g, '_')}_export.xlsx`);
  };

  // Flatten hierarchy for export
  const flattenHierarchy = (data: RowData[]): any[] => {
    const result: any[] = [];
    const flatten = (rows: RowData[]) => {
      rows.forEach(row => {
        result.push(row);
        if (row.children) flatten(row.children);
      });
    };
    flatten(data);
    return result;
  };

  // Handle add comment (stub from your original)
  const handleAddComment = () => {
    if (!commentCell || !newComment.trim()) return;
    const { rowId, colKey } = commentCell;
    const comment = { user: 'Me', text: newComment, time: new Date().toLocaleString() };

    setCellComments(prev => {
      const key = `${rowId}-${colKey}`;
      return { ...prev, [key]: [...(prev[key] || []), comment] };
    });

    // Backend sync stub
    console.log('Syncing comment to backend:', { rowId, colKey, comment });
    setNewComment('');
  };

  // ... (rest of the component: renderTreeRows, pivot logic, panels, etc. unchanged)

  return (
    <DndProvider backend={HTML5Backend}>
      <div ref={gridRef} className="flex flex-col h-full p-4 bg-gray-950 text-gray-100 rounded-xl border border-gray-800">
        {/* Header */}
        <div className="flex justify-between items-center mb-4 flex-wrap gap-4">
          <h2 className="text-xl font-semibold">{appId}</h2>
          <div className="flex items-center gap-3 flex-wrap">
            <button
              onClick={() => setGridStyle('vector')}
              className={`px-3 py-1.5 rounded-md border ${gridStyle === 'vector' ? 'bg-indigo-600 border-indigo-500' : 'bg-gray-900 border-gray-700'}`}
            >
              Vector
            </button>
            <button
              onClick={() => setGridStyle('matrix')}
              className={`px-3 py-1.5 rounded-md border ${gridStyle === 'matrix' ? 'bg-indigo-600 border-indigo-500' : 'bg-gray-900 border-gray-700'}`}
            >
              Matrix
            </button>
            {/* ... all previous buttons: undo/redo, collab, pivot, validation, chart, export */}
          </div>
        </div>

        {/* Dimension Slicers */}
        <div className="flex flex-wrap gap-4 mb-6 p-4 bg-gray-900/70 rounded-lg border border-gray-800">
          {/* ... unchanged */}
        </div>

        {/* Loading */}
        {loading && (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
          </div>
        )}

        {/* Main Grid */}
        {!loading && gridStyle === 'vector' && (
          <div className="flex-1 overflow-hidden">
            <div className="overflow-auto h-full">
              <table className="w-full text-left border-collapse relative">
                {/* ... thead, tbody, renderTreeRows unchanged */}
              </table>
            </div>
          </div>
        )}

        {!loading && gridStyle === 'matrix' && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <div className="rounded-lg border border-gray-800 bg-gray-900/60 p-3">
              <h3 className="text-sm font-semibold mb-2">Yield Input (Month x Year)</h3>
              <table className="w-full text-xs">
                <thead>
                  <tr>
                    <th className="text-left py-1">Month</th>
                    {MATRIX_YEARS.map((year) => (
                      <th key={year} className="py-1">{year}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {MATRIX_MONTHS.map((month, monthIndex) => (
                    <tr key={month}>
                      <td className="py-1">{month}</td>
                      {MATRIX_YEARS.map((year) => (
                        <td key={year} className="py-1">
                          <input
                            type="number"
                            value={String(yieldMatrix[monthIndex][year] ?? 0)}
                            onChange={(event) => handleMatrixInput(setYieldMatrix, monthIndex, year, event.target.value)}
                            className="w-16 bg-gray-950 border border-gray-700 rounded px-1 py-0.5"
                            aria-label={`Yield ${month} ${year}`}
                            title={`Yield ${month} ${year}`}
                            placeholder="0"
                          />
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="rounded-lg border border-gray-800 bg-gray-900/60 p-3">
              <h3 className="text-sm font-semibold mb-2">Area Input (Month x Year)</h3>
              <table className="w-full text-xs">
                <thead>
                  <tr>
                    <th className="text-left py-1">Month</th>
                    {MATRIX_YEARS.map((year) => (
                      <th key={year} className="py-1">{year}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {MATRIX_MONTHS.map((month, monthIndex) => (
                    <tr key={month}>
                      <td className="py-1">{month}</td>
                      {MATRIX_YEARS.map((year) => (
                        <td key={year} className="py-1">
                          <input
                            type="number"
                            value={String(areaMatrix[monthIndex][year] ?? 0)}
                            onChange={(event) => handleMatrixInput(setAreaMatrix, monthIndex, year, event.target.value)}
                            className="w-16 bg-gray-950 border border-gray-700 rounded px-1 py-0.5"
                            aria-label={`Area ${month} ${year}`}
                            title={`Area ${month} ${year}`}
                            placeholder="0"
                          />
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="rounded-lg border border-gray-800 bg-gray-900/60 p-3">
              <h3 className="text-sm font-semibold mb-2">Weight = Yield × Area</h3>
              <table className="w-full text-xs">
                <thead>
                  <tr>
                    <th className="text-left py-1">Month</th>
                    {MATRIX_YEARS.map((year) => (
                      <th key={year} className="py-1">{year}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {MATRIX_MONTHS.map((month, monthIndex) => (
                    <tr key={month}>
                      <td className="py-1">{month}</td>
                      {MATRIX_YEARS.map((year) => (
                        <td key={year} className="py-1 text-indigo-300 font-medium">{Number(weightMatrix[monthIndex][year] ?? 0).toFixed(2)}</td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="rounded-lg border border-gray-800 bg-gray-900/60 p-3">
              <h3 className="text-sm font-semibold mb-2">Revenue = Price × Weight</h3>
              <table className="w-full text-xs">
                <thead>
                  <tr>
                    <th className="text-left py-1">Month</th>
                    {MATRIX_YEARS.map((year) => (
                      <th key={year} className="py-1">{year}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {MATRIX_MONTHS.map((month, monthIndex) => (
                    <tr key={month}>
                      <td className="py-1">{month}</td>
                      {MATRIX_YEARS.map((year) => (
                        <td key={year} className="py-1 text-emerald-300 font-medium">{Number(revenueMatrix[monthIndex][year] ?? 0).toFixed(2)}</td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="rounded-lg border border-gray-800 bg-gray-900/60 p-3 lg:col-span-2">
              <h3 className="text-sm font-semibold mb-2">Vector Projection (staging conversion)</h3>
              <p className="text-xs text-gray-400 mb-2">
                Matrix input is converted to vector rows by explicit period qualifier (<span className="font-mono">year-month</span>) before revenue math.
              </p>
              <div className="overflow-auto max-h-64">
                <table className="w-full text-xs">
                  <thead>
                    <tr>
                      <th className="text-left py-1 pr-2">Period</th>
                      <th className="py-1 pr-2">Yield</th>
                      <th className="py-1 pr-2">Area</th>
                      <th className="py-1 pr-2">Weight</th>
                      <th className="py-1 pr-2">Price</th>
                      <th className="py-1 pr-2">Revenue</th>
                    </tr>
                  </thead>
                  <tbody>
                    {vectorProjectionRows.map((row) => (
                      <tr key={row.periodKey} className="border-t border-gray-800">
                        <td className="py-1 pr-2">{row.month} {row.year}</td>
                        <td className="py-1 pr-2 text-center">{row.yield.toFixed(2)}</td>
                        <td className="py-1 pr-2 text-center">{row.area.toFixed(2)}</td>
                        <td className="py-1 pr-2 text-center text-indigo-300">{row.weight.toFixed(2)}</td>
                        <td className="py-1 pr-2 text-center">{row.price.toFixed(2)}</td>
                        <td className="py-1 pr-2 text-center text-emerald-300">{row.revenue.toFixed(2)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        )}

        {/* Pivot, Chart, Panels unchanged */}
      </div>
    </DndProvider>
  );
}

// ────────────────────────────────────────────────
// Dashboard (unchanged)
// ────────────────────────────────────────────────

export function PlanningDashboard() {
  const salesDimensions: Dimension[] = [
    { name: 'Time', members: ['2025', '2026-Q1', '2026-Q2', '2026-Q3', '2026-Q4', 'All'] },
    { name: 'Entity', members: ['North America', 'EMEA', 'APAC', 'All'] },
    { name: 'Product Line', members: ['Hardware', 'Software', 'Services', 'All'] }
  ];

  const hrDimensions: Dimension[] = [
    { name: 'Department', members: ['Engineering', 'Sales', 'Finance', 'HR', 'All'] },
    { name: 'Scenario', members: ['Budget', 'Forecast', 'Actual', 'All'] },
    { name: 'Year', members: ['2025', '2026', 'All'] }
  ];

  return (
    <div className="min-h-screen bg-gray-950 p-6 space-y-12">
      <h1 className="text-3xl font-bold text-center text-indigo-400">Planning & Reporting Dashboard</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="bg-gray-900/50 rounded-xl border border-gray-800 p-4">
          <QuantAtomGrid
            appId="Sales & Revenue Planning"
            cubeEndpoint="/api/olap/sales"
            dimensions={salesDimensions}
            columnOrder={['entity', 'time', 'product', 'budget', 'actual', 'variance']}
          />
        </div>

        <div className="bg-gray-900/50 rounded-xl border border-gray-800 p-4">
          <QuantAtomGrid
            appId="Headcount & Compensation"
            cubeEndpoint="/api/olap/hr"
            dimensions={hrDimensions}
            columnOrder={['department', 'headcount', 'salary', 'bonus', 'total']}
          />
        </div>
      </div>
    </div>
  );
}