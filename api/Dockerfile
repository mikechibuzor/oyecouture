# Go API — multi-stage build producing a small static image.
# Migrations are embedded in the binary via go:embed, so no SQL files are copied.
FROM golang:1.23-alpine AS build

WORKDIR /src

# Dependencies first so this layer caches between source-only changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static binary: no libc dependency, so it runs on the distroless base below.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/server \
    ./cmd/server

# Distroless: no shell, no package manager — a much smaller attack surface than alpine.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/server /server

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/server"]
