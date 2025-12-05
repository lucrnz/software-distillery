FROM alpine/curl:latest AS curl-version
RUN curl --version | awk -F' ' '{print $2}' > /curl-version.txt

FROM golang:1.25.5-alpine AS builder
ARG COMMIT_HASH=unknown
ENV COMMIT_HASH=${COMMIT_HASH}
ENV GOOS=linux
WORKDIR /app
COPY . .

COPY --from=curl-version /curl-version.txt /curl-version.txt
RUN CURL_VERSION="$(cat /curl-version.txt)" && \
  rm -f /curl-version.txt && \
  apk add --no-cache make && \
  make build

FROM scratch
COPY --from=builder /app/build/simple-downloader /simple-downloader

ENTRYPOINT ["/simple-downloader"]
