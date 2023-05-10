### builder ###
FROM golang:1.19 as builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
  -ldflags "-s -w" -tags osusergo,netgo -a -o emtectl ./cmd/emtectl

### runner ###
FROM alpine:3.17

LABEL org.opencontainers.image.authors="Shota Kitazawa, Kohei Ota"
LABEL org.opencontainers.image.url="https://github.com/cloudnativedaysjp/emtec-ecu"
LABEL org.opencontainers.image.source="https://github.com/cloudnativedaysjp/emtec-ecu/blob/main/cmd/server/Dockerfile"
WORKDIR /
COPY --from=builder /workspace/emtectl .

ENTRYPOINT ["/emtectl"]
