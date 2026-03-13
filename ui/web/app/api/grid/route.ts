import { NextResponse } from "next/server";

const GRID_SERVICE_URL = process.env.GRID_SERVICE_URL || "http://localhost:8080/grid/query";

const DEFAULT_PAYLOAD = {
  dimensions: {
    rows: ["Entity"],
    columns: ["Time"],
    pages: [],
    filters: {},
  },
  members: {
    Entity: ["North America", "EMEA", "APAC"],
    Time: ["2025", "2026", "2027"],
    Measure: ["Revenue"],
    Scenario: ["Actual"],
  },
  window: { rowStart: 0, rowEnd: 0, colStart: 0, colEnd: 0 },
  defaults: {},
  stream: false,
};

export async function POST(req: Request) {
  const body = await req.json().catch(() => ({}));
  const payload = Object.keys(body || {}).length ? body : DEFAULT_PAYLOAD;

  try {
    const upstream = await fetch(GRID_SERVICE_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!upstream.ok) {
      const detail = await upstream.text();
      console.warn("/api/grid: upstream error", upstream.status, detail);
      return NextResponse.json(mockGrid(detail), { status: 200 });
    }

    const json = await upstream.json();
    return NextResponse.json(json);
  } catch (err: any) {
    console.warn("/api/grid: upstream unreachable", err?.message || err);
    return NextResponse.json(mockGrid(err?.message || String(err)), { status: 200 });
  }
}

function mockGrid(reason: string) {
  const rows = [
    [{ id: 1, code: "NA", name: "North America" }],
    [{ id: 2, code: "EMEA", name: "EMEA" }],
    [{ id: 3, code: "APAC", name: "APAC" }],
  ];
  const columns = [
    [{ id: 101, code: "2025", name: "2025" }],
    [{ id: 102, code: "2026", name: "2026" }],
    [{ id: 103, code: "2027", name: "2027" }],
  ];
  const cells = rows.flatMap((_, r) =>
    columns.map((_, c) => ({ rowIndex: r, colIndex: c, value: 1000 * (r + 1) + 100 * c }))
  );
  return { rows, columns, cells, fallback: true, reason } as any;
}
