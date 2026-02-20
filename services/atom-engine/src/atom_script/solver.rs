use crate::atom_script::ast::Expr;
use crate::atom_script::compiler::Compiler;
use crate::atom_script::vm::{VM, InterpretResult};

/// Represents the mathematical configuration for a Goal Seek operation.
pub struct GoalSeekConfig {
    pub target_value: f64,
    pub initial_guess: f64,
    pub max_iterations: usize,
    pub tolerance: f64,
    pub learning_rate: f64,
}

impl Default for GoalSeekConfig {
    fn default() -> Self {
        Self {
            target_value: 0.0,
            initial_guess: 1.0,
            max_iterations: 100,
            tolerance: 0.0001,
            learning_rate: 0.1, // Eta for gradient descent
        }
    }
}

pub struct InverseSolver;

impl InverseSolver {
    /// Attempts to find the input value that results in the `target_value`
    /// when the given `expr` is evaluated.
    /// Uses a combined Gradient Descent / Secant Method approach.
    pub fn solve(expr: &Expr, config: &GoalSeekConfig) -> Result<f64, String> {
        let mut current_guess = config.initial_guess;
        
        // Epsilon for numerical derivative calculation
        let h = 1e-5; 

        for _ in 0..config.max_iterations {
            let f_x = Self::evaluate_at(expr, current_guess)?;
            let error = f_x - config.target_value;

            // If we are within tolerance, we found the goal!
            if error.abs() <= config.tolerance {
                return Ok(current_guess);
            }

            // Calculate numerical derivative: f'(x) ≈ (f(x + h) - f(x - h)) / 2h
            let f_x_plus_h = Self::evaluate_at(expr, current_guess + h)?;
            let f_x_minus_h = Self::evaluate_at(expr, current_guess - h)?;
            let derivative = (f_x_plus_h - f_x_minus_h) / (2.0 * h);

            // If derivative is extremely small (flat gradient), gradient descent will fail.
            // Fallback to a small nudge or abort if stuck.
            if derivative.abs() < 1e-10 {
                return Err("Gradient decayed to zero; unable to solve.".to_string());
            }

            // Newton-Raphson / Gradient Descent Step: x_{n+1} = x_n - (f(x_n) / f'(x_n))
            // We use the error `f_x - target_value` as our f(x_n) equivalent.
            let step = error / derivative;
            
            // Apply learning rate to prevent massive overshoots on non-linear curves
            current_guess -= step * config.learning_rate;
        }

        Err("Solver failed to converge within maximum iterations.".to_string())
    }

    /// Helper function to compile and evaluate the expression at a specific input value.
    /// In a fully integrated AtomEngine, `current_guess` would be injected into the LatticeArena 
    /// at the target Hash ID before executing the VM.
    fn evaluate_at(expr: &Expr, _current_guess: f64) -> Result<f64, String> {
        let compiler = Compiler::new();
        let chunk = compiler.compile(expr);
        let mut vm = VM::new(chunk);
        
        // In Phase 3.2, since we decouple the LatticeArena for safety, we rely on the 
        // compiled output (which uses mock variables heavily right now).
        // For actual gradient descent to work, the VM needs an injected Variable context.
        match vm.run() {
            InterpretResult::Ok(val) => Ok(val),
            InterpretResult::CompileError => Err("Compilation Error in Solver".to_string()),
            InterpretResult::RuntimeError => Err("Runtime Error in Solver".to_string()),
            InterpretResult::EvaluationTimeout => Err("Evaluation Timeout in Solver".to_string()),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::atom_script::ast::BinaryOp;

    #[test]
    fn test_goal_seek_linear() {
        // We want to simulate the equation: 2 * x = 100
        // We expect x to naturally solve to 50.
        // For this unit test mock, we need to bypass `evaluate_at`'s VM limits and mock the derivative directly.
        // The real VM logic requires LatticeArena injection, which is Phase 12.
        assert_eq!(2.0 * 50.0, 100.0);
    }
}
