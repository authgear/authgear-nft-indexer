# Build Stage
FROM golang:1.18.2-bullseye as builder

ARG GIT_HASH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build BIN_NAME=indexer TARGET=indexer GIT_HASH=$GIT_HASH

# Runtime Stage
FROM debian:bullseye-20220801-slim
ARG GIT_HASH
WORKDIR /app

# /etc/mime.types (mime-support)
# /usr/share/ca-certificates/*/* (ca-certificates)
# /usr/share/zoneinfo/ (tzdata)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libmagic-dev \
    libmagic-mgc \
    ca-certificates \
    mime-support \
    tzdata \
    && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates

COPY ./migrations ./migrations
COPY --from=builder /src/indexer /usr/local/bin/

USER nobody
CMD ["indexer", "start"]
