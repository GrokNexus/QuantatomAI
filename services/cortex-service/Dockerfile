# Build Stage
FROM python:3.12-slim as builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip wheel --no-cache-dir --no-deps --wheel-dir /app/wheels -r requirements.txt

# Final Stage
FROM python:3.12-slim

WORKDIR /app

# Copy built wheels
COPY --from=builder /app/wheels /wheels
COPY --from=builder /app/requirements.txt .

# Install dependencies from wheels
RUN pip install --no-cache /wheels/*

# Copy application source
COPY src/ ./src/

# Expose FastAPI port
EXPOSE 8081

# Run the Uvicorn server
CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8081", "--workers", "4"]
