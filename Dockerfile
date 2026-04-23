## Build
FROM  --platform=$BUILDPLATFORM golang:1.26 AS server
WORKDIR /app

# Set build arguments for target OS and architecture.
# These will be automatically populated by 'docker buildx' based on the requested platform.
ARG TARGETOS
ARG TARGETARCH

# Set up environment variables for Go cross-compilation
# CGO_ENABLED=0 creates a statically linked binary, which is essential for the scratch stage.
ENV CGO_ENABLED=0
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

COPY go.mod go.sum ./
RUN go mod download

COPY --parents **/*.go ./

RUN GOEXPERIMENT=jsonv2 go build -trimpath -ldflags="-w -s" -o main ./cmd/main.go

FROM --platform=$BUILDPLATFORM node:25-bookworm AS frontend

WORKDIR /app

COPY ui/package.json ui/package-lock.json ./
RUN npm ci

COPY ui .
RUN npm run build

## Deploy
FROM alpine:3.20
WORKDIR /
COPY --from=server /app/main /usr/bin/
COPY --from=frontend /app/dist static

ENV STATIC_DIR=/static
ENTRYPOINT ["main"]