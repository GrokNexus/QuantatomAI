use std::collections::HashMap;

use thiserror::Error;

#[derive(Debug, Clone, PartialEq)]
pub enum Comparator {
    Eq,
    NotEq,
    Gt,
    Gte,
    Lt,
    Lte,
}

#[derive(Debug, Clone, PartialEq)]
pub enum FilterValue {
    Str(String),
    Num(f64),
    Ref(String),
}

#[derive(Debug, Clone, PartialEq)]
pub struct FilterPredicate {
    pub left_ref: String,
    pub comparator: Comparator,
    pub right: FilterValue,
}

#[derive(Debug, Clone, PartialEq)]
pub enum Expr {
    FunctionCall {
        name: String,
        args: Vec<String>,
        filter: Option<FilterPredicate>,
    },
}

#[derive(Debug, Error)]
pub enum ParseError {
    #[error("invalid formula: {0}")]
    Invalid(String),
}

#[derive(Debug, Error)]
pub enum EvalError {
    #[error("unknown function: {0}")]
    UnknownFunction(String),
    #[error("missing metric: {0}")]
    MissingMetric(String),
}

#[derive(Debug, Clone)]
pub struct Record {
    pub metrics: HashMap<String, f64>,
    pub attrs: HashMap<String, String>,
}

impl Record {
    pub fn new() -> Self {
        Self {
            metrics: HashMap::new(),
            attrs: HashMap::new(),
        }
    }

    pub fn with_metric(mut self, key: &str, value: f64) -> Self {
        self.metrics.insert(key.to_string(), value);
        self
    }

    pub fn with_attr(mut self, key: &str, value: &str) -> Self {
        self.attrs.insert(key.to_string(), value.to_string());
        self
    }
}

pub fn parse_formula(input: &str) -> Result<Expr, ParseError> {
    let trimmed = input.trim();
    if trimmed.is_empty() {
        return Err(ParseError::Invalid("empty input".to_string()));
    }

    let upper = trimmed.to_uppercase();
    let where_index = upper.find(" WHERE ");

    let (call_part, where_part) = match where_index {
        Some(idx) => (&trimmed[..idx], Some(trimmed[idx + 7..].trim())),
        None => (trimmed, None),
    };

    let open_paren = call_part
        .find('(')
        .ok_or_else(|| ParseError::Invalid("missing '('".to_string()))?;
    let close_paren = call_part
        .rfind(')')
        .ok_or_else(|| ParseError::Invalid("missing ')'".to_string()))?;

    if close_paren <= open_paren {
        return Err(ParseError::Invalid("malformed function call".to_string()));
    }

    let function_name = call_part[..open_paren].trim();
    if function_name.is_empty() {
        return Err(ParseError::Invalid("missing function name".to_string()));
    }

    let raw_args = call_part[open_paren + 1..close_paren].trim();
    let args = if raw_args.is_empty() {
        Vec::new()
    } else {
        raw_args
            .split(',')
            .map(|s| normalize_arg(s.trim()))
            .collect::<Result<Vec<_>, _>>()?
    };

    let filter = match where_part {
        Some(expr) => Some(parse_predicate(expr)?),
        None => None,
    };

    Ok(Expr::FunctionCall {
        name: function_name.to_uppercase(),
        args,
        filter,
    })
}

fn normalize_arg(raw: &str) -> Result<String, ParseError> {
    if raw.starts_with('[') && raw.ends_with(']') && raw.len() > 2 {
        Ok(raw[1..raw.len() - 1].to_string())
    } else if !raw.is_empty() {
        Ok(raw.to_string())
    } else {
        Err(ParseError::Invalid("empty argument".to_string()))
    }
}

fn parse_predicate(expr: &str) -> Result<FilterPredicate, ParseError> {
    let candidates = [
        ("!=", Comparator::NotEq),
        (">=", Comparator::Gte),
        ("<=", Comparator::Lte),
        ("=", Comparator::Eq),
        (">", Comparator::Gt),
        ("<", Comparator::Lt),
    ];

    for (token, cmp) in candidates {
        if let Some(idx) = expr.find(token) {
            let left = expr[..idx].trim();
            let right = expr[idx + token.len()..].trim();

            if left.starts_with('[') && left.ends_with(']') && left.len() > 2 {
                let left_ref = left[1..left.len() - 1].to_string();
                let right_value = parse_filter_value(right)?;
                return Ok(FilterPredicate {
                    left_ref,
                    comparator: cmp,
                    right: right_value,
                });
            }

            return Err(ParseError::Invalid("left predicate must be [Reference]".to_string()));
        }
    }

    Err(ParseError::Invalid("missing predicate comparator".to_string()))
}

fn parse_filter_value(raw: &str) -> Result<FilterValue, ParseError> {
    if raw.starts_with('[') && raw.ends_with(']') && raw.len() > 2 {
        return Ok(FilterValue::Ref(raw[1..raw.len() - 1].to_string()));
    }

    if raw.starts_with('"') && raw.ends_with('"') && raw.len() >= 2 {
        return Ok(FilterValue::Str(raw[1..raw.len() - 1].to_string()));
    }

    if let Ok(value) = raw.parse::<f64>() {
        return Ok(FilterValue::Num(value));
    }

    Err(ParseError::Invalid(format!("invalid predicate value: {raw}")))
}

pub fn evaluate(expr: &Expr, records: &[Record]) -> Result<f64, EvalError> {
    match expr {
        Expr::FunctionCall { name, args, filter } => {
            if name != "SUM" {
                return Err(EvalError::UnknownFunction(name.clone()));
            }

            if args.is_empty() {
                return Ok(0.0);
            }

            let metric = &args[0];
            let filtered = records.iter().filter(|r| passes_filter(r, filter));

            let mut sum = 0.0;
            for record in filtered {
                let value = record
                    .metrics
                    .get(metric)
                    .ok_or_else(|| EvalError::MissingMetric(metric.clone()))?;
                sum += value;
            }

            Ok(sum)
        }
    }
}

fn passes_filter(record: &Record, filter: &Option<FilterPredicate>) -> bool {
    let Some(predicate) = filter else {
        return true;
    };

    let left_val = record.attrs.get(&predicate.left_ref).cloned().unwrap_or_default();

    match &predicate.right {
        FilterValue::Str(right) => compare_str(&left_val, right, &predicate.comparator),
        FilterValue::Ref(right_ref) => {
            let right_val = record.attrs.get(right_ref).cloned().unwrap_or_default();
            compare_str(&left_val, &right_val, &predicate.comparator)
        }
        FilterValue::Num(right) => {
            let parsed_left = left_val.parse::<f64>().unwrap_or(0.0);
            compare_num(parsed_left, *right, &predicate.comparator)
        }
    }
}

fn compare_str(left: &str, right: &str, cmp: &Comparator) -> bool {
    match cmp {
        Comparator::Eq => left == right,
        Comparator::NotEq => left != right,
        Comparator::Gt => left > right,
        Comparator::Gte => left >= right,
        Comparator::Lt => left < right,
        Comparator::Lte => left <= right,
    }
}

fn compare_num(left: f64, right: f64, cmp: &Comparator) -> bool {
    match cmp {
        Comparator::Eq => (left - right).abs() < f64::EPSILON,
        Comparator::NotEq => (left - right).abs() >= f64::EPSILON,
        Comparator::Gt => left > right,
        Comparator::Gte => left >= right,
        Comparator::Lt => left < right,
        Comparator::Lte => left <= right,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parses_sum_reference() {
        let expr = parse_formula("SUM([Revenue])").unwrap();
        match expr {
            Expr::FunctionCall { name, args, filter } => {
                assert_eq!(name, "SUM");
                assert_eq!(args, vec!["Revenue"]);
                assert!(filter.is_none());
            }
        }
    }

    #[test]
    fn parses_sum_where_not_equal() {
        let expr = parse_formula("SUM([Revenue]) WHERE [Region] != \"Intercompany\"").unwrap();
        match expr {
            Expr::FunctionCall { filter, .. } => {
                let p = filter.expect("filter expected");
                assert_eq!(p.left_ref, "Region");
                assert_eq!(p.comparator, Comparator::NotEq);
                assert_eq!(p.right, FilterValue::Str("Intercompany".to_string()));
            }
        }
    }

    #[test]
    fn parses_eq_predicate() {
        let expr = parse_formula("SUM([COGS]) WHERE [Scenario] = \"Budget\"").unwrap();
        match expr {
            Expr::FunctionCall { filter, .. } => {
                let p = filter.expect("filter expected");
                assert_eq!(p.comparator, Comparator::Eq);
            }
        }
    }

    #[test]
    fn parses_gte_numeric_predicate() {
        let expr = parse_formula("SUM([Revenue]) WHERE [Score] >= 90").unwrap();
        match expr {
            Expr::FunctionCall { filter, .. } => {
                let p = filter.expect("filter expected");
                assert_eq!(p.comparator, Comparator::Gte);
                assert_eq!(p.right, FilterValue::Num(90.0));
            }
        }
    }

    #[test]
    fn parses_lt_numeric_predicate() {
        let expr = parse_formula("SUM([Revenue]) WHERE [Score] < 90").unwrap();
        match expr {
            Expr::FunctionCall { filter, .. } => {
                let p = filter.expect("filter expected");
                assert_eq!(p.comparator, Comparator::Lt);
                assert_eq!(p.right, FilterValue::Num(90.0));
            }
        }
    }

    #[test]
    fn parses_ref_to_ref_predicate() {
        let expr = parse_formula("SUM([Revenue]) WHERE [Region] = [ReportingRegion]").unwrap();
        match expr {
            Expr::FunctionCall { filter, .. } => {
                let p = filter.expect("filter expected");
                assert_eq!(p.right, FilterValue::Ref("ReportingRegion".to_string()));
            }
        }
    }

    #[test]
    fn rejects_missing_function_name() {
        let err = parse_formula("([Revenue])").unwrap_err();
        assert!(format!("{err}").contains("missing function name"));
    }

    #[test]
    fn rejects_missing_parenthesis() {
        let err = parse_formula("SUM[Revenue]").unwrap_err();
        assert!(format!("{err}").contains("missing '('") || format!("{err}").contains("missing ')'"));
    }

    #[test]
    fn evaluates_sum_without_filter() {
        let expr = parse_formula("SUM([Revenue])").unwrap();
        let rows = vec![
            Record::new().with_metric("Revenue", 10.0),
            Record::new().with_metric("Revenue", 20.0),
            Record::new().with_metric("Revenue", 5.5),
        ];

        let got = evaluate(&expr, &rows).unwrap();
        assert!((got - 35.5).abs() < 0.0001);
    }

    #[test]
    fn evaluates_sum_with_string_filter() {
        let expr = parse_formula("SUM([Revenue]) WHERE [Region] != \"Intercompany\"").unwrap();
        let rows = vec![
            Record::new()
                .with_metric("Revenue", 10.0)
                .with_attr("Region", "NA"),
            Record::new()
                .with_metric("Revenue", 20.0)
                .with_attr("Region", "Intercompany"),
            Record::new()
                .with_metric("Revenue", 5.0)
                .with_attr("Region", "EMEA"),
        ];

        let got = evaluate(&expr, &rows).unwrap();
        assert!((got - 15.0).abs() < 0.0001);
    }
}
