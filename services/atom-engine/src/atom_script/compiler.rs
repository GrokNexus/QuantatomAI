use crate::atom_script::ast::{BinaryOp, Expr};
use crate::atom_script::chunk::{Chunk, OpCode};
use crate::lattice::metadata::{HierarchyResolver, MockHierarchyResolver};

pub struct Compiler {
    chunk: Chunk,
    resolver: Box<dyn HierarchyResolver>,
}

impl Compiler {
    pub fn new() -> Self {
        Self {
            chunk: Chunk::new(),
            resolver: Box::new(MockHierarchyResolver), // Default to Mock for now
        }
    }

    pub fn compile(mut self, expr: &Expr) -> Chunk {
        self.compile_expr(expr);
        self.chunk.write_chunk(OpCode::Return);
        self.chunk
    }

    fn compile_expr(&mut self, expr: &Expr) {
        self.compile_expr_with_count(expr);
    }

    /// Compiles an expression and returns the number of values pushed to the stack.
    /// Usually 1, but can be N for Hierarchy Expansions.
    fn compile_expr_with_count(&mut self, expr: &Expr) -> usize {
        match expr {
            Expr::Literal(val) => {
                let idx = self.chunk.add_constant(*val);
                self.chunk.write_chunk(OpCode::Constant(idx));
                1
            }
            Expr::Binary { op, lhs, rhs } => {
                // Optimization: Constant Folding
                if let (Expr::Literal(l), Expr::Literal(r)) = (lhs.as_ref(), rhs.as_ref()) {
                     let val = match op {
                         BinaryOp::Add => l + r,
                         BinaryOp::Sub => l - r,
                         BinaryOp::Mul => l * r,
                         BinaryOp::Div => l / r,
                     };
                     let idx = self.chunk.add_constant(val);
                     self.chunk.write_chunk(OpCode::Constant(idx));
                     return 1;
                }

                self.compile_expr(lhs);
                self.compile_expr(rhs);
                match op {
                    BinaryOp::Add => self.chunk.write_chunk(OpCode::Add),
                    BinaryOp::Sub => self.chunk.write_chunk(OpCode::Sub),
                    BinaryOp::Mul => self.chunk.write_chunk(OpCode::Mul),
                    BinaryOp::Div => self.chunk.write_chunk(OpCode::Div),
                }
                1
            }
            Expr::Identifier(_) => {
                // TODO: Load variable
                1
            }
            Expr::DimensionRef(name) => {
                // TODO: Emit OpCode::LoadDimension(name)
                // For now, push 0.0 placeholder
                let idx = self.chunk.add_constant(0.0);
                self.chunk.write_chunk(OpCode::Constant(idx));
                1
            }
            Expr::FunctionCall { name, args } => {
                let mut arg_count = 0;
                for arg in args {
                    arg_count += self.compile_expr_with_count(arg);
                }
                
                match name.as_str() {
                    "SUM" => self.chunk.write_chunk(OpCode::Sum(arg_count)),
                    "AVG" => self.chunk.write_chunk(OpCode::Avg(arg_count)),
                    "MIN" => self.chunk.write_chunk(OpCode::Min(arg_count)),
                    "MAX" => self.chunk.write_chunk(OpCode::Max(arg_count)),
                    "LOOKUP" => self.chunk.write_chunk(OpCode::Lookup),
                    "XLOOKUP" => self.chunk.write_chunk(OpCode::XLookup(arg_count)),
                    _ => {
                        // TODO: Unknown function
                    }
                }
                1 
            }
            // Ultra Diamond: Hierarchy Expansion
            Expr::HierarchyCall { name, args } => {
                if name == "Children" && args.len() == 2 {
                    if let (Expr::DimensionRef(dim), Expr::DimensionRef(member)) = (&args[0], &args[1]) {
                        let children = self.resolver.get_children(dim, member);
                        let count = children.len();
                        for child in children {
                             // Emit Load for each child
                             // Re-using compile logic for DimensionRef
                             self.compile_expr_with_count(&Expr::DimensionRef(child));
                        }
                        return count;
                    }
                }
                0 // Error or empty
            }
            // Ultra Diamond: Time Travel (Shift)
            Expr::TimeTravel { lhs, rhs } => {
                self.compile_expr(lhs);
                self.compile_expr(rhs);
                self.chunk.write_chunk(OpCode::Shift);
                1
            }
        }
    }
}
