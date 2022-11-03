### builder ###
FROM golang:1.19 as builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags osusergo,netgo -a -o emtectl ./cmd/emtectl

### runner ###
FROM alpine:3.16.2
WORKDIR /
COPY --from=builder /workspace/emtectl .

ENTRYPOINT ["/emtectl"]
