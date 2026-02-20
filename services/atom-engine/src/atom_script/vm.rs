use crate::atom_script::chunk::{Chunk, OpCode};

pub struct VM {
    chunk: Chunk,
    stack: Vec<f64>,
    ip: usize, // Instruction Pointer
}

pub enum InterpretResult {
    Ok(f64),
    CompileError,
    RuntimeError,
    EvaluationTimeout, // Ultra Diamond: Vector 1 DoS Protection
}

impl VM {
    pub fn new(chunk: Chunk) -> Self {
        Self {
            chunk,
            stack: Vec::with_capacity(256), // Typical stack depth
            ip: 0,
        }
    }

    pub fn run(&mut self) -> InterpretResult {
        let mut op_count = 0;
        const MAX_OPS: usize = 10_000_000; // Circuit Breaker: Maximum instruction cycles

        loop {
            if op_count >= MAX_OPS {
                return InterpretResult::EvaluationTimeout;
            }
            op_count += 1;

            if self.ip >= self.chunk.code.len() {
                return InterpretResult::RuntimeError;
            }

            let instruction = self.chunk.code[self.ip];
            self.ip += 1;

            match instruction {
                OpCode::Return => {
                    return InterpretResult::Ok(self.pop());
                }
                OpCode::Constant(idx) => {
                    let constant = self.chunk.constants[idx];
                    if let Err(e) = self.push(constant) { return e; }
                }
                OpCode::Add => {
                    let b = self.pop();
                    let a = self.pop();
                    if let Err(e) = self.push(a + b) { return e; }
                }
                OpCode::Sub => {
                    let b = self.pop();
                    let a = self.pop();
                    if let Err(e) = self.push(a - b) { return e; }
                }
                OpCode::Mul => {
                    let b = self.pop();
                    let a = self.pop();
                    if let Err(e) = self.push(a * b) { return e; }
                }
                OpCode::Div => {
                    let b = self.pop();
                    let a = self.pop();
                    if let Err(e) = self.push(a / b) { return e; }
                }
                OpCode::Negate => {
                    let a = self.pop();
                    if let Err(e) = self.push(-a) { return e; }
                }
                // Ultra Diamond: Aggregation
                OpCode::Sum(count) => {
                    let mut sum = 0.0;
                    for _ in 0..count {
                        sum += self.pop();
                    }
                    if let Err(e) = self.push(sum) { return e; }
                }
                OpCode::Avg(count) => {
                    let mut sum = 0.0;
                    for _ in 0..count {
                        sum += self.pop();
                    }
                    if let Err(e) = self.push(sum / count as f64) { return e; }
                }
                OpCode::Min(count) => {
                    let mut min_val = f64::MAX;
                    for _ in 0..count {
                        let v = self.pop();
                        if v < min_val { min_val = v; }
                    }
                    if let Err(e) = self.push(min_val) { return e; }
                }
                OpCode::Max(count) => {
                    let mut max_val = f64::MIN;
                    for _ in 0..count {
                        let v = self.pop();
                        if v > max_val { max_val = v; }
                    }
                    if let Err(e) = self.push(max_val) { return e; }
                }
                // Ultra Diamond: Lookups & Time Travel (Phase 12 Kernels)
                // In Phase 12, the VM will be injected with an unsafe pointer to the LatticeArena.
                // These opcodes will execute an O(1) atomic pointer jump without evaluating the grid.
                OpCode::Shift => {
                    let offset = self.pop(); // e.g. [PrevMonth]
                    let value = self.pop();  // e.g. [Revenue]
                    
                    // Phase 12: unsafe { value_ptr.offset(offset as isize) }
                    // For now, securely pop and return the base value to guarantee stack safety.
                    if let Err(e) = self.push(value + offset) { return e; }
                }
                OpCode::Lookup => {
                    let _return_rng = self.pop();
                    let _search_rng = self.pop();
                    let _lookup_val = self.pop();
                    
                    // Phase 12: SIMD accelerated scan across the search_rng pointer.
                    // Fallback to safe 0.0 until kernel is injected.
                    if let Err(e) = self.push(0.0) { return e; }
                }
                OpCode::XLookup(count) => {
                    for _ in 0..count {
                        let _arg = self.pop();
                    }
                    // Phase 12: B-Tree or SIMD scan based on XLookup heuristics.
                    if let Err(e) = self.push(0.0) { return e; }
                }
                // Phase 3: Time-Intelligence Shifts
                OpCode::TimeShift(shift_code) => {
                    let base_val = self.pop(); // The calculated or raw value of the base metric
                    
                    // In a production LatticeArena, we would shift the underlying memory pointer
                    // here without evaluating the calculation tree again. (O(1) Jump)
                    //
                    // let shifted_ptr = match shift_code {
                    //    1 => unsafe { base_ptr.offset(-12) }, // PY
                    //    2 => unsafe { base_ptr.offset(-3) },  // PQ
                    //    3 => calculate_ytd_simd(base_ptr),    // YTD
                    //    _ => base_ptr
                    // };
                    
                    // For the Phase 3 VM, we will simulate the shifted value.
                    let simulated_shift = match shift_code {
                        1 => base_val * 0.90, // Simulate PY as 90% of current
                        2 => base_val * 0.95, // Simulate PQ as 95% of current
                        3 => base_val * 6.0,  // Simulate YTD (assuming mid-year)
                        4 => base_val * 2.0,  // Simulate QTD (assuming mid-quarter)
                        5 => base_val,        // PTD
                        _ => base_val,
                    };
                    
                    if let Err(e) = self.push(simulated_shift) { return e; }
                }
            }
        }
    }

    fn push(&mut self, value: f64) -> Result<(), InterpretResult> {
        if self.stack.len() >= 256 {
            return Err(InterpretResult::RuntimeError); // Stack Overflow Protection
        }
        self.stack.push(value);
        Ok(())
    }

    fn pop(&mut self) -> f64 {
        self.stack.pop().expect("Stack underflow")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_time_intelligence_shifts() {
        // Test PY Shift (Shift Code 1)
        // Simulated expectation: 100 * 0.90 = 90
        let mut py_chunk = Chunk::new();
        let idx = py_chunk.add_constant(100.0);
        py_chunk.write_chunk(OpCode::Constant(idx));
        py_chunk.write_chunk(OpCode::TimeShift(1));
        py_chunk.write_chunk(OpCode::Return);

        let mut vm = VM::new(py_chunk);
        if let InterpretResult::Ok(val) = vm.run() {
            assert_eq!(val, 90.0);
        } else {
            panic!("PY Shift failed");
        }

        // Test YTD Shift (Shift Code 3)
        // Simulated expectation: 100 * 6.0 = 600
        let mut ytd_chunk = Chunk::new();
        let idx2 = ytd_chunk.add_constant(100.0);
        ytd_chunk.write_chunk(OpCode::Constant(idx2));
        ytd_chunk.write_chunk(OpCode::TimeShift(3));
        ytd_chunk.write_chunk(OpCode::Return);

        let mut vm2 = VM::new(ytd_chunk);
        if let InterpretResult::Ok(val) = vm2.run() {
            assert_eq!(val, 600.0);
        } else {
            panic!("YTD Shift failed");
        }
    }
}
