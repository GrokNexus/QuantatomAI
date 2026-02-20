
use arrow::datatypes::{DataType, Field, Schema};
use std::sync::Arc;

pub struct MoleculeSchema;

impl MoleculeSchema {
    pub fn schema() -> Arc<Schema> {
        Arc::new(Schema::new(vec![
            Field::new("coordinate_hash", DataType::Binary, false),
            // Map is complex in Arrow, often represented as List of Structs. 
            // For simplicity in v1, we might treat custom_dimensions as List<Struct<Key, Value>>
            // But Parquet Go writer writes Map. Arrow reader should handle it.
            // Let's rely on schema inference from file for now.
            // Field::new("custom_dimensions", DataType::Map(...), true), 
            
            Field::new("numeric_value", DataType::Float64, true),
            Field::new("text_commentary", DataType::Utf8, true),
            Field::new("embedding_vector", DataType::Binary, true),
            
            // Ultra Diamond Upgrade: Rich Types
            Field::new("date_value", DataType::Int64, true), // Unix Millis
            Field::new("boolean_value", DataType::Boolean, true),
            Field::new("error_value", DataType::Utf8, true),
            
            Field::new("timestamp", DataType::Int64, false),
            Field::new("source_system", DataType::Utf8, false),
            Field::new("security_mask", DataType::UInt64, false),
            Field::new("causality_clock", DataType::Binary, true),
            
            // Collision Resolution (Red Team Audit Response)
            Field::new("is_locked", DataType::Boolean, true),
        ]))
    }
}
