use crate::atom_script::parser::Parser;
use crate::atom_script::compiler::Compiler;
use crate::atom_script::chunk::OpCode;

#[test]
fn test_hierarchy_children_expansion() {
    // 1. Parsing
    // @Children expands to members: USA, Canada, Mexico
    let input = "SUM(@Children([Region], [North America]))";
    let mut parser = Parser::new(input);
    let expr = parser.parse().expect("Parse failed");

    // 2. Compilation
    let compiler = Compiler::new();
    let chunk = compiler.compile(&expr);

    // 3. Verification
    // We expect:
    // - 3 Constants (placeholders for USA, Canada, Mexico)
    // - OpCode::Sum(3)
    
    // Check if OpCode::Sum(3) is present
    let has_sum_3 = chunk.code.iter().any(|op| *op == OpCode::Sum(3));
    assert!(has_sum_3, "Chunk should contain OpCode::Sum(3). Code: {:?}", chunk.code);

    // Check if we have constants loaded. 
    // Since DimensionRef currently pushes a 0.0 constant, we should have at least 3 constants.
    assert!(chunk.constants.len() >= 3, "Should have at least 3 constants for the children");
}

#[test]
fn test_lookup_and_time_travel() {
    // 1. Parsing LOOKUP
    let input = "LOOKUP(10, [Range], [Return])";
    let mut parser = Parser::new(input);
    let expr = parser.parse().expect("Parse failed");
    let compiler = Compiler::new();
    let chunk = compiler.compile(&expr);
    assert!(chunk.code.contains(&OpCode::Lookup));

    // 2. Parsing Time Travel
    let input = "[Revenue] -> [PrevMonth]";
    let mut parser = Parser::new(input);
    let expr = parser.parse().expect("Parse failed");
    let compiler = Compiler::new();
    let chunk = compiler.compile(&expr);
    assert!(chunk.code.contains(&OpCode::Shift));
}

#[test]
fn test_basic_arithmetic() {
    let input = "1 + 2";
    let mut parser = Parser::new(input);
    let expr = parser.parse().expect("Parse failed");
    
    let compiler = Compiler::new();
    let chunk = compiler.compile(&expr);

    // Constant Folding should reduce this to a single constant (3.0)
    assert!(chunk.code.contains(&OpCode::Constant(0))); 
    assert_eq!(chunk.constants[0], 3.0);
}
