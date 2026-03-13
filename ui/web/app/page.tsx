"use client";

import { Grid } from "../grid-engine/core/Grid";

export default function Home() {
    return (
        <main className="p-6 space-y-4 bg-slate-950 min-h-screen">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-semibold text-slate-100">Grid Engine (Live)</h1>
                    <p className="text-slate-400 text-sm">Served from grid-service via /api/grid proxy.</p>
                </div>
            </div>
            <Grid />
        </main>
    );
}