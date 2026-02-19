#[derive(Debug, Clone, Copy, PartialEq)]
pub enum OpCode {
    Return,
    Constant(usize), // Index in constants pool
    Add,
    Sub,
    Mul,
    Div,
    Negate,
    // Ultra Diamond: Aggregation Ops
    Sum(usize), // Pops N items from stack
    Avg(usize),
    Min(usize),
    Max(usize),

    // Ultra Diamond: Lookups & Time Travel
    Lookup, // Pops 3: range, search_val, return_range
    XLookup(usize), // Pops N (Standard args)
    Shift, // Pops 2: Dimension, Offset/Target
}

pub struct Chunk {
    pub code: Vec<OpCode>,
    pub constants: Vec<f64>,
}

impl Chunk {
    pub fn new() -> Self {
        Self {
            code: Vec::new(),
            constants: Vec::new(),
        }
    }

    pub fn write_chunk(&mut self, byte: OpCode) {
        self.code.push(byte);
    }

    pub fn add_constant(&mut self, value: f64) -> usize {
        self.constants.push(value);
        self.constants.len() - 1
    }
}
