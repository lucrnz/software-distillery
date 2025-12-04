FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=0
ENV GOOS=linux

ADD . .
WORKDIR /app

RUN go build -ldflags="-s -w" -o simple-downloader main.go

FROM stratch AS target
COPY --from=builder /app/simple-downloader /simple-downloader

ENTRYPOINT [ "/simple-downloader" ]
