use logos::Logos;

#[derive(Logos, Debug, PartialEq)]
pub enum Token {
    // Arithmetic Operators
    #[token("+")]
    Plus,
    #[token("-")]
    Minus,
    #[token("*")]
    Mul,
    #[token("/")]
    Div,
    #[token("(")]
    LParen,
    #[token(")")]
    RParen,
    #[token(",")]
    Comma,

    // Excel-style Functions
    #[token("SUM")]
    Sum,
    #[token("AVG")]
    Avg,
    #[token("MIN")]
    Min,
    #[token("MAX")]
    Max,
    #[token("IF")]
    If,
    #[token("LOOKUP")]
    Lookup,
    #[token("XLOOKUP")]
    XLookup,

    // Time Travel Operator
    #[token("->")]
    Arrow,

    // Identifiers (e.g., [Region], Revenue)
    #[regex("[a-zA-Z_][a-zA-Z0-9_]*", |lex| lex.slice().to_string())]
    Identifier(String),

    // Hierarchy Functions (e.g., @Children)
    #[regex("@[a-zA-Z_][a-zA-Z0-9_]*", |lex| lex.slice()[1..].to_string())]
    AtIdentifier(String),

    // Dimension References (e.g., [Region])
    #[regex(r"\[[^\]]*\]", |lex| lex.slice().trim_matches(|c| c == '[' || c == ']').to_string())]
    DimensionRef(String),

    // Number Literals
    #[regex(r"[0-9]+(\.[0-9]+)?", |lex| lex.slice().parse().ok())]
    Number(f64),

    // Whitespace
    #[regex(r"[ \t\n\f]+", logos::skip)]
    Whitespace, // Ignored

    // Logos 0.13+ error handling
    Error,
}
