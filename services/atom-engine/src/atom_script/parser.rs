use logos::{Logos, Lexer};
use crate::atom_script::lexer::Token;
use crate::atom_script::ast::{Expr, BinaryOp};

pub struct Parser<'a> {
    lexer: Lexer<'a, Token>,
    current_token: Option<Token>,
}

impl<'a> Parser<'a> {
    pub fn new(input: &'a str) -> Self {
        let mut lexer = Token::lexer(input);
        let first_token = lexer.next().map(|res| res.unwrap_or(Token::Error));
        Self {
            lexer,
            current_token: first_token,
        }
    }

    fn advance(&mut self) {
        self.current_token = self.lexer.next().map(|res| res.unwrap_or(Token::Error));
    }

    pub fn parse(&mut self) -> Result<Expr, String> {
        self.parse_expr(0)
    }

    fn parse_expr(&mut self, min_bp: u8) -> Result<Expr, String> {
        let mut lhs = match &self.current_token {
            Some(Token::Number(n)) => {
                let val = *n;
                self.advance();
                Expr::Literal(val)
            }
            Some(Token::DimensionRef(d)) => {
                let name = d.clone();
                self.advance();
                Expr::DimensionRef(name)
            }
            Some(Token::Identifier(id)) => {
                let name = id.clone();
                self.advance();
                // Check if function call
                if let Some(Token::LParen) = self.current_token {
                    self.advance();
                    // parse args
                    let args = self.parse_args()?;
                    Expr::FunctionCall { name, args }
                } else {
                    Expr::Identifier(name)
                }
            }
            // Ultra Diamond: Hierarchy Functions (@Children)
            Some(Token::AtIdentifier(id)) => {
                let name = id.clone();
                self.advance();
                if self.current_token != Some(Token::LParen) {
                    return Err("Expected '(' after hierarchy function".to_string());
                }
                self.advance();
                let args = self.parse_args()?;
                Expr::HierarchyCall { name, args }
            }
            Some(Token::Lookup) => {
                self.advance();
                if self.current_token != Some(Token::LParen) { return Err("Expected '(' after LOOKUP".to_string()); }
                self.advance();
                let args = self.parse_args()?;
                Expr::FunctionCall { name: "LOOKUP".to_string(), args }
            }
            Some(Token::XLookup) => {
                self.advance();
                if self.current_token != Some(Token::LParen) { return Err("Expected '(' after XLOOKUP".to_string()); }
                self.advance();
                let args = self.parse_args()?;
                Expr::FunctionCall { name: "XLOOKUP".to_string(), args }
            }
            Some(Token::Sum) => {
                self.advance();
                if self.current_token != Some(Token::LParen) { return Err("Expected '(' after SUM".to_string()); }
                self.advance();
                let args = self.parse_args()?;
                Expr::FunctionCall { name: "SUM".to_string(), args }
            }
            Some(Token::Avg) => {
                self.advance();
                if self.current_token != Some(Token::LParen) { return Err("Expected '(' after AVG".to_string()); }
                self.advance();
                let args = self.parse_args()?;
                Expr::FunctionCall { name: "AVG".to_string(), args }
            }
            Some(Token::Min) => {
                self.advance();
                if self.current_token != Some(Token::LParen) { return Err("Expected '(' after MIN".to_string()); }
                self.advance();
                let args = self.parse_args()?;
                Expr::FunctionCall { name: "MIN".to_string(), args }
            }
            Some(Token::Max) => {
                self.advance();
                if self.current_token != Some(Token::LParen) { return Err("Expected '(' after MAX".to_string()); }
                self.advance();
                let args = self.parse_args()?;
                Expr::FunctionCall { name: "MAX".to_string(), args }
            }
            Some(Token::LParen) => {
                self.advance();
                let expr = self.parse_expr(0)?;
                if self.current_token != Some(Token::RParen) {
                    return Err("Expected ')'".to_string());
                }
                self.advance();
                expr
            }
            _ => return Err(format!("Unexpected token: {:?}", self.current_token)),
        };

        loop {
            // Ultra Diamond: Time Travel Operator (->)
            if let Some(Token::Arrow) = &self.current_token {
                let (l_bp, r_bp) = (5, 6); // High precedence
                if l_bp < min_bp { break; }
                self.advance();
                let rhs = self.parse_expr(r_bp)?;
                lhs = Expr::TimeTravel { lhs: Box::new(lhs), rhs: Box::new(rhs) };
                continue;
            }

            let op = match &self.current_token {
                Some(Token::Plus) => BinaryOp::Add,
                Some(Token::Minus) => BinaryOp::Sub,
                Some(Token::Mul) => BinaryOp::Mul,
                Some(Token::Div) => BinaryOp::Div,
                _ => break,
            };

            let (l_bp, r_bp) = infix_binding_power(&op);
            if l_bp < min_bp {
                break;
            }

            self.advance();
            let rhs = self.parse_expr(r_bp)?;
            lhs = Expr::Binary {
                op,
                lhs: Box::new(lhs),
                rhs: Box::new(rhs),
            };
        }

        Ok(lhs)
    }

    fn parse_args(&mut self) -> Result<Vec<Expr>, String> {
        let mut args = Vec::new();
        if self.current_token != Some(Token::RParen) {
            loop {
                args.push(self.parse_expr(0)?);
                if self.current_token == Some(Token::Comma) {
                    self.advance();
                } else {
                    break;
                }
            }
        }
        if self.current_token != Some(Token::RParen) {
                return Err("Expected ')'".to_string());
        }
        self.advance();
        Ok(args)
    }
}

fn infix_binding_power(op: &BinaryOp) -> (u8, u8) {
    match op {
        BinaryOp::Add | BinaryOp::Sub => (1, 2),
        BinaryOp::Mul | BinaryOp::Div => (3, 4),
    }
}
