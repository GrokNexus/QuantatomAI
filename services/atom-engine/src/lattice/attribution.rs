use std::cmp::Ordering;
use crate::lattice::metadata::HierarchyResolver;
use crate::lattice::arena::LatticeArena;
use rayon::prelude::*;

#[derive(Debug, Clone, PartialEq)]
pub struct VarianceDriver {
    pub dimension: String,
    pub member: String,
    pub scenario_a_val: f64,
    pub scenario_b_val: f64,
    pub variance: f64,
    pub percentage_of_total: f64,
}

pub struct AttributionEngine;

impl AttributionEngine {
    /// Analyzes a top-level variance and traverses the HierarchyResolver to find the Top N 
    /// statistical drivers that caused the variance, pulling live data from the LatticeArena.
    pub fn calculate_drivers<'a, R: HierarchyResolver>(
        resolver: &'a R,
        arena: &LatticeArena,
        dimension: &str,
        parent_member: &str,
        _scenario_a: &str,
        _scenario_b: &str,
        top_n: usize,
        base_hash_a: u128, // Base hash for Scenario A intersection
        base_hash_b: u128, // Base hash for Scenario B intersection
    ) -> Vec<VarianceDriver> {
        
        // 1. Get all immediate children from the metadata hierarchy
        let children = resolver.get_children(dimension, parent_member);
        
        // 2. Query the LatticeArena for each child concurrently (Ultra-Diamond Rayon Parallelism)
        let mut drivers: Vec<VarianceDriver> = children.par_iter().map(|child| {
            // In a real system, the hash is a combination of the base hash and the child's ID.
            // For Phase 8.2 proof-of-concept, we simulate the specific hashes.
            let hash_a = base_hash_a.wrapping_add(child.len() as u128); // dummy hash combinator
            let hash_b = base_hash_b.wrapping_add(child.len() as u128);

            let val_a = arena.get_cell(hash_a);
            let val_b = arena.get_cell(hash_b);
            let variance = val_b - val_a;

            VarianceDriver {
                dimension: dimension.to_string(),
                member: child.clone(),
                scenario_a_val: val_a,
                scenario_b_val: val_b,
                variance,
                percentage_of_total: 0.0, // calculated later
            }
        }).collect();

        let total_absolute_variance: f64 = drivers.iter().map(|d| d.variance.abs()).sum();

        // 3. Compute Percentage of Total Abs Variance (to feed the LLM later)
        for driver in &mut drivers {
            if total_absolute_variance != 0.0 {
                driver.percentage_of_total = driver.variance.abs() / total_absolute_variance;
            }
        }

        // 4. Sort by Absolute Variance Descending
        drivers.sort_by(|a, b| {
            b.variance.abs().partial_cmp(&a.variance.abs()).unwrap_or(Ordering::Equal)
        });

        // 5. Return the Top N Drivers to be sent to the Python Fluxion Engine
        drivers.into_iter().take(top_n).collect()
    }
}

// Tests
#[cfg(test)]
mod tests {
    use super::*;
    use crate::lattice::metadata::MockHierarchyResolver;

    #[test]
    fn test_attribution_engine() {
        let arena = LatticeArena::new(1000);
        let resolver = MockHierarchyResolver;

        // Seed some data for the children of "North America" ("USA"-3, "Canada"-6, "Mexico"-6)
        let hash_a_usa = 100u128.wrapping_add(3);
        let hash_b_usa = 200u128.wrapping_add(3);
        arena.set_cell(hash_a_usa, 100.0);
        arena.set_cell(hash_b_usa, 150.0); // +50 variance

        let hash_a_can = 100u128.wrapping_add(6);
        let hash_b_can = 200u128.wrapping_add(6);
        arena.set_cell(hash_a_can, 50.0);
        arena.set_cell(hash_b_can, 40.0); // -10 variance

        let drivers = AttributionEngine::calculate_drivers(
            &resolver,
            &arena,
            "Region",
            "North America",
            "Actual_2024",
            "Actual_2025",
            2,
            100, // base_a
            200, // base_b
        );

        assert_eq!(drivers.len(), 2);
        assert_eq!(drivers[0].member, "USA"); // Highest variance (50)
        assert_eq!(drivers[0].variance, 50.0);
        assert_eq!(drivers[1].member, "Canada"); // Next highest variance (-10)
    }
}
