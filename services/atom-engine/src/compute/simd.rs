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

    /// Spreads a `target` value proportionally across cells based on `reference_values`.
    /// Respects the `is_locked` bitmask to prevent overwriting explicit bottom-up entries.
    /// The remaining target is spread across the unlocked cells.
    pub fn proportional_spread(
        target: f64, 
        current_values: &[f64], 
        reference_values: &[f64], 
        is_locked: &[bool]
    ) -> Vec<f64> {
        // Step 1: Calculate the total locked value that has already been spoken for.
        let locked_sum: f64 = current_values.par_iter()
            .zip(is_locked.par_iter())
            .filter_map(|(&val, &locked)| if locked { Some(val) } else { None })
            .sum();

        let remaining_target = target - locked_sum;

        // Step 2: Calculate the total reference weight of the UNLOCKED cells.
        let unlocked_ref_sum: f64 = reference_values.par_iter()
            .zip(is_locked.par_iter())
            .filter_map(|(&ref_val, &locked)| if !locked { Some(ref_val) } else { None })
            .sum();

        // Avoid divide by zero if all unlocked reference cells sum to 0
        let safe_ref_sum = if unlocked_ref_sum == 0.0 { 1.0 } else { unlocked_ref_sum };

        // Step 3: Compute the new values in a single wait-free parallel pass.
        current_values.par_iter()
            .zip(reference_values.par_iter())
            .zip(is_locked.par_iter())
            .map(|((&cur, &ref_val), &locked)| {
                if locked {
                    cur // Keep the explicit bottom-up entry
                } else if unlocked_ref_sum == 0.0 {
                    0.0 // Could distribute evenly here, but default to 0 to prevent Sparsity Explosion
                } else {
                    (ref_val / safe_ref_sum) * remaining_target
                }
            })
            .collect()
    }
}
