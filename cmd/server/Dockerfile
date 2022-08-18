# Build Stage
FROM golang:1.18.2-bullseye as builder

ARG GIT_HASH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build BIN_NAME=server TARGET=server GIT_HASH=$GIT_HASH

# Runtime Stage
FROM debian:bullseye-20220801-slim
ARG GIT_HASH
WORKDIR /app

COPY ./migrations ./migrations
COPY --from=builder /src/server /usr/local/bin/

USER nobody
CMD ["server", "start"]