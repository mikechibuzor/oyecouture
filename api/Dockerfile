# Go API with a multi-stage build producing a small static image.
# Migrations are embedded in the binary via go:embed, so no SQL files are copied.
# Must match the go directive in go.mod: GOTOOLCHAIN is local in this image, so
# an older base fails the build rather than downloading a newer toolchain.
FROM golang:1.25-alpine AS build

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

# The admin CLI ships alongside the server: without it there is no way to seed
# the owner account on a deployed box, and the admin is unreachable.
#   docker compose exec -it api /admin
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/admin \
    ./cmd/admin

# Distroless has no shell or package manager, giving it a smaller attack surface than alpine.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/server /server
COPY --from=build /out/admin /admin

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/server"]
