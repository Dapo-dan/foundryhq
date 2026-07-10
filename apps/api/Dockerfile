# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder
WORKDIR /src

# Copy module files first so `go mod download` is cached separately from
# source changes — avoids re-downloading deps on every code edit.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Passed via --build-arg (e.g. VERSION=$(git describe --tags), COMMIT=$(git
# rev-parse --short HEAD)); defaults match internal/version's own fallback
# so an unparameterized build still works.
ARG VERSION=dev
ARG COMMIT=unknown
# CGO_ENABLED=0 produces a static binary (the postgres driver is pure Go via
# pgx, so no CGO is needed) that can run on a bare alpine base below.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "\
    -X github.com/foundryhq/foundryhq/apps/api/internal/version.Version=${VERSION} \
    -X github.com/foundryhq/foundryhq/apps/api/internal/version.Commit=${COMMIT}" \
    -o /out/server ./cmd/server

FROM alpine:3.20
# ca-certificates is needed for any outbound TLS calls (e.g. OAuth providers
# once those ship); harmless to include now.
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /out/server ./server

EXPOSE 8080
ENTRYPOINT ["./server"]
