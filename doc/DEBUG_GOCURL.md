## Debug

You might want to add [gocurl (A Go implementation of curl)](https://github.com/ameshkov/gocurl) to the created image to test stuff (example: upload files to temp file hosting, etc),

As often old containers might have outdated/insecure versions of curl.

For this add these steps:

```Dockerfile
# Download gocurl
FROM curl AS gocurl-download
RUN GOCURL_VER="1.5.0" && \
  URL_PREFIX="https://github.com/ameshkov/gocurl/releases/download/v${GOCURL_VER}/" && \
  ARCH=$(uname -m) && \
  case "$ARCH" in \
      x86_64) \
          FILE="gocurl-linux-amd64-v${GOCURL_VER}.tar.gz"; \
          PLATFORM="amd64"; \
          HASH="70e5cfe5fb0cc3468a538d55fc493e1f44392462d86d183340c16029be7aa779" ;; \
      aarch64) \
          FILE="gocurl-linux-arm64-v${GOCURL_VER}.tar.gz"; \
          PLATFORM="arm64"; \
          HASH="9363a2e393ff360db20399e3fe37f6600bad7611aa78dd747952dedadd042dd5" ;; \
      armv7l|armv7*) \
          FILE="gocurl-linux-arm-v${GOCURL_VER}.tar.gz"; \
          PLATFORM="arm"; \
          HASH="797db2f7ba3d14cb4b29f987cf3ce9815e2ba6835adfda2a89a2c72aacda9cb8" ;; \
      i686|i386) \
          FILE="gocurl-linux-386-v${GOCURL_VER}.tar.gz"; \
          PLATFORM="386"; \
          HASH="9d324eb8a1a38de63386fa094d3f9d858209a8006aa535d5f7b5282c2ca53df6" ;; \
      *) echo "Unsupported architecture: $ARCH"; exit 1 ;; \
  esac && \
  curl -fLO "${URL_PREFIX}${FILE}" && \
  echo "$HASH $FILE" | sha256sum -c - && \
  tar -xzf "$FILE" --strip-components=1 -C /usr/local/bin/ "linux-${PLATFORM}/gocurl" && \
  chmod +x /usr/local/bin/gocurl && \
  rm "$FILE"
```

On the target image stage, add almost at the end:

```Dockerfile
# Install gocurl
COPY --from=gocurl-download /usr/local/bin/gocurl /usr/local/bin/gocurl
RUN chmod +x /usr/local/bin/gocurl
```
