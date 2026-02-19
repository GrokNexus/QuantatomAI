use std::fs::File;
use parquet::arrow::arrow_reader::ParquetRecordBatchReaderBuilder;
use arrow::record_batch::RecordBatch;
use anyhow::Result;

/// Reads an MDF (Parquet) file into a vector of Arrow RecordBatches.
/// This uses Zero-Copy semantics where possible, mapping the file directly into memory.
pub fn read_mdf_arrow(path: &str) -> Result<Vec<RecordBatch>> {
    let file = File::open(path)?;
    
    // Create the builder from the file
    let builder = ParquetRecordBatchReaderBuilder::try_new(file)?;
    
    // Build the reader with optimal batch size for SIMD processing (e.g., 8192 rows)
    let reader = builder.with_batch_size(8192).build()?;
    
    // Collect all batches (streams into memory)
    let batches: Result<Vec<_>, _> = reader.collect();
    
    Ok(batches?)
}
