FROM golang:1.25.5-alpine AS builder
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app
ADD . .

RUN go build -ldflags="-s -w" -o simple-downloader

FROM scratch
COPY --from=builder /app/simple-downloader /simple-downloader

ENTRYPOINT [ "/simple-downloader" ]
