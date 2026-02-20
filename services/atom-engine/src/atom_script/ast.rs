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
        rhs: Box<Expr>, // previously [PrevMonth], now we transition to TimeModifier
    },
    // Phase 3: Core Time Modifiers (PY, PQ, YTD, etc)
    TimeModifier {
        base: Box<Expr>, // The metric to shift: [Revenue]
        shift_type: TimeShiftType,
    },
}

#[derive(Debug, PartialEq, Clone)]
pub enum TimeShiftType {
    PriorYear,
    PriorQuarter,
    YearToDate,
    QuarterToDate,
    PeriodToDate,
}
