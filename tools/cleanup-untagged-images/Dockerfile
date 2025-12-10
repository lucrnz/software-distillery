# syntax=docker/dockerfile:1
FROM python:3.13-alpine AS builder

# Install uv
COPY --from=ghcr.io/astral-sh/uv:0.9.17 /uv /usr/local/bin/uv

# Set working directory
WORKDIR /build

# Copy project files
COPY pyproject.toml ./
COPY cleanup.py ./

# Create virtual environment and install project
RUN set -euxo pipefail && \
    uv venv /opt/venv && \
    uv pip install --python /opt/venv/bin/python .

# Final stage
FROM python:3.13-alpine

# Copy virtual environment from builder
COPY --from=builder /opt/venv /opt/venv
COPY --from=builder /build/cleanup.py /app/cleanup.py

# Set environment to use virtual environment
ENV PATH="/opt/venv/bin:${PATH}"

# Set working directory
WORKDIR /app

# Run the cleanup script
ENTRYPOINT ["python", "/app/cleanup.py"]
