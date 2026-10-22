# syntax=docker/dockerfile:1

# --- Build stage ---
FROM golang:1.27-alpine AS build
WORKDIR /src

# No third-party dependencies, so go.mod alone is enough to prime the module cache.
COPY go.mod ./
RUN go mod download

COPY . .
# Static binary (no cgo) so it runs on a distroless image without libc.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# --- Runtime stage ---
# Distroless: no shell or package manager, ~2 MB base, and a built-in non-root user.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/server"]
