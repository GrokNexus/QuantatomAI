#[derive(Debug, PartialEq, Clone)]
pub enum BinaryOp {
    Add,
    Sub,
    Mul,
    Div,
}

#[derive(Debug, PartialEq, Clone)]
pub enum Expr {
    Literal(f64),
    Identifier(String),
    DimensionRef(String), // e.g. [Region]
    Binary {
        op: BinaryOp,
        lhs: Box<Expr>,
        rhs: Box<Expr>,
    },
    FunctionCall {
        name: String,
        args: Vec<Expr>,
    },
    // Ultra Diamond: Hierarchy Macro
    HierarchyCall {
        name: String,
        args: Vec<Expr>,
    },
    // Ultra Diamond: Time Travel
    TimeTravel {
        lhs: Box<Expr>, // e.g. [Revenue]
        rhs: Box<Expr>, // e.g. [PrevMonth]
    },
}
