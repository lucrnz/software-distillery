FROM ghcr.io/astral-sh/uv:0.9-python3.14-alpine AS builder
WORKDIR /build

COPY pyproject.toml ./
COPY cleanup.py ./

RUN set -euxo pipefail && \
  uv venv /opt/venv && \
  uv pip install --python /opt/venv/bin/python .

FROM python:3.14-alpine

COPY --from=builder /opt/venv /opt/venv
COPY --from=builder /build/cleanup.py /app/cleanup.py
ENV PATH="/opt/venv/bin:${PATH}"

WORKDIR /app

ENTRYPOINT ["python", "/app/cleanup.py"]
