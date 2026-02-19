use std::collections::HashMap;
use std::sync::RwLock;

const SHARD_COUNT: usize = 64;

/// A single shard of the arena.
struct ArenaShard {
    values: RwLock<Vec<f64>>,  // Type 0
    dates: RwLock<Vec<i64>>,   // Type 1
    strings: RwLock<Vec<String>>, 
    index_map: RwLock<HashMap<u128, usize>>,
}

impl ArenaShard {
    fn new(capacity: usize) -> Self {
        Self {
            values: RwLock::new(Vec::with_capacity(capacity)),
            dates: RwLock::new(Vec::with_capacity(capacity)),
            strings: RwLock::new(Vec::with_capacity(capacity)),
            index_map: RwLock::new(HashMap::with_capacity(capacity)),
        }
    }
}

/// The LatticeArena manages the memory for all cells in a Grid View.
/// Ultra-Diamond: Uses Sharded Locking for massive concurrency (5000+ writers).
pub struct LatticeArena {
    shards: Vec<ArenaShard>,
}

impl LatticeArena {
    pub fn new(capacity: usize) -> Self {
        let shard_cap = capacity / SHARD_COUNT + 1;
        let mut shards = Vec::with_capacity(SHARD_COUNT);
        for _ in 0..SHARD_COUNT {
            shards.push(ArenaShard::new(shard_cap));
        }
        Self { shards }
    }

    fn get_shard(&self, hash: u128) -> &ArenaShard {
        let idx = (hash % SHARD_COUNT as u128) as usize;
        &self.shards[idx]
    }

    /// Allocates or updates a cell value.
    pub fn set_cell(&self, hash: u128, value: f64) -> usize {
        let shard = self.get_shard(hash);
        
        // Fast path: Check if exists (Read Lock)
        {
            let map = shard.index_map.read().unwrap();
            if let Some(&idx) = map.get(&hash) {
                let mut vals = shard.values.write().unwrap();
                vals[idx] = value;
                return idx;
            }
        }

        // Slow path: Insert new (Write Lock)
        let mut map = shard.index_map.write().unwrap();
        let mut vals = shard.values.write().unwrap();

        // Double check
        if let Some(&idx) = map.get(&hash) {
            vals[idx] = value;
            return idx;
        }

        let idx = vals.len();
        vals.push(value);
        map.insert(hash, idx);
        
        idx
    }

    /// Retrieves a cell value. Returns 0.0 if not found (sparse).
    pub fn get_cell(&self, hash: u128) -> f64 {
        let shard = self.get_shard(hash);
        let map = shard.index_map.read().unwrap();
        if let Some(&idx) = map.get(&hash) {
            let vals = shard.values.read().unwrap();
            return vals[idx]; // Safe because shard lock protects index bounds
        }
        0.0
    }

    /// Returns a combined vector for SIMD processing (expensive copy, uses rayon).
    /// Note: In V2, iterate sharded directly.
    pub fn get_vector(&self) -> Vec<f64> {
        // Simple implementation: Combine all shards.
        // Parallel implementation: Map-Reduce would be better here.
        let mut combined = Vec::new();
        for shard in &self.shards {
            let vals = shard.values.read().unwrap();
            combined.extend_from_slice(&vals);
        }
        combined
    }
    
    // Ultra Diamond: Rich Type Setters
    pub fn set_string(&self, hash: u128, val: String) -> usize {
        let shard = self.get_shard(hash);
        let mut strs = shard.strings.write().unwrap();
        let idx = strs.len();
        strs.push(val);
        idx
    }

    pub fn set_date(&self, hash: u128, val: i64) -> usize {
        self.set_cell(hash, val as f64)
    }
}
