import React, { useState } from 'react';
import { GridCanvas } from './components/GridCanvas';
import { FormulaBar } from './components/FormulaBar';

function App() {
    const [formula, setFormula] = useState("Sum([Revenue])");

    return (
        <div className="App" style={{ height: '100vh', width: '100vw', backgroundColor: '#1e1e1e', overflow: 'hidden' }}>
            <h1 style={{ color: 'white', position: 'absolute', top: '50px', zIndex: 10, padding: '10px' }}>
                QuantatomAI Grid (WebGPU)
            </h1>

            {/* Layer 7.3: Formula Editor */}
            <div style={{ position: 'absolute', top: 0, left: 0, width: '100%', zIndex: 100, backgroundColor: '#252526' }}>
                <FormulaBar value={formula} onChange={setFormula} />
            </div>

            <div style={{ paddingTop: '45px', width: '100%', height: '100%' }}>
                <GridCanvas />
            </div>
        </div>
    );
}

export default App;
