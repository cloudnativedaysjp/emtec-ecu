### builder ###
FROM golang:1.19 as builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags osusergo,netgo -a -o server ./cmd/server

### runner ###
FROM gcr.io/distroless/static-debian11:nonroot
WORKDIR /
COPY --from=builder /workspace/server .

ENTRYPOINT ["/cndctl"]
