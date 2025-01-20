# Build Stage
FROM quay.io/theauthgear/golang:1.23.4-noble as builder

ARG GIT_HASH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build BIN_NAME=server TARGET=server GIT_HASH=$GIT_HASH

# Runtime Stage
FROM ubuntu:noble
ARG GIT_HASH
WORKDIR /app

# /etc/mime.types (media-types)
# /usr/share/ca-certificates/*/* (ca-certificates)
# /usr/share/zoneinfo/ (tzdata)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libmagic-dev \
    libmagic-mgc \
    ca-certificates \
    media-types \
    tzdata \
    && apt-get upgrade -y && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates

COPY --from=builder /src/server /usr/local/bin/

USER nobody
CMD ["server", "start"]
