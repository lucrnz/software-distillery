FROM ghcr.io/astral-sh/uv:0.9-python3.14-alpine AS builder
WORKDIR /build

COPY pyproject.toml ./
COPY ghcr_cleanup ./ghcr_cleanup

RUN set -euxo pipefail && \
  uv venv /opt/venv && \
  uv pip install --python /opt/venv/bin/python .

FROM python:3.14-alpine

COPY --from=builder /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:${PATH}"

WORKDIR /app

ENTRYPOINT ["ghcr-cleanup"]
