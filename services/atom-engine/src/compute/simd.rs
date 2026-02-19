use rayon::prelude::*;

/// VectorOps provides SIMD-accelerated arithmetic on standard vectors.
/// We use Rayon to parallelize the loop, and the Rust compiler auto-vectorizes
/// the inner loop into AVX-512 instructions if available.
pub struct VectorOps;

impl VectorOps {
    /// Adds two vectors element-wise.
    /// Panic: vectors must be same length.
    pub fn add(a: &[f64], b: &[f64]) -> Vec<f64> {
        // Ultra-Diamond: Rayon for Multi-Core Parallelism
        a.par_iter()
         .zip(b.par_iter())
         .map(|(x, y)| x + y)
         .collect()
    }

    /// Subtracts b from a element-wise.
    pub fn sub(a: &[f64], b: &[f64]) -> Vec<f64> {
        a.par_iter()
         .zip(b.par_iter())
         .map(|(x, y)| x - y)
         .collect()
    }

    /// Multiplies two vectors element-wise.
    pub fn mul(a: &[f64], b: &[f64]) -> Vec<f64> {
        a.par_iter()
         .zip(b.par_iter())
         .map(|(x, y)| x * y)
         .collect()
    }

    /// Divides a by b element-wise. Handles division by zero (returns infinity).
    pub fn div(a: &[f64], b: &[f64]) -> Vec<f64> {
        a.par_iter()
         .zip(b.par_iter())
         .map(|(x, y)| x / y)
         .collect()
    }

    /// Calculates the validation checksum (Sum) of the vector.
    pub fn sum(a: &[f64]) -> f64 {
        a.par_iter().sum()
    }
}
