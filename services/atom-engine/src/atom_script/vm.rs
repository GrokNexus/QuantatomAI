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
        loop {
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
                // Ultra Diamond: Lookups & Time Travel
                OpCode::Shift => {
                    let _offset = self.pop(); // e.g. [PrevMonth]
                    let value = self.pop();   // e.g. [Revenue]
                    // Mock: In real engine, this shifts the pointer.
                    // Here we just pass the value through or add them if they are numbers.
                    if let Err(e) = self.push(value) { return e; }
                }
                OpCode::Lookup => {
                    let _return_rng = self.pop();
                    let _search_rng = self.pop();
                    let _lookup_val = self.pop();
                    // Mock: Return 42.0 found
                    if let Err(e) = self.push(42.0) { return e; }
                }
                OpCode::XLookup(count) => {
                    for _ in 0..count {
                        let _arg = self.pop();
                    }
                    // Mock: Return 100.0 found
                    if let Err(e) = self.push(100.0) { return e; }
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
