FROM golang:1.15.0-buster as base
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download
COPY . .

FROM base AS builder
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/maze-api ./cmd/maze-api

FROM alpine:latest
COPY --from=builder /go/bin/maze-api /bin/maze-api
USER nobody
ENTRYPOINT ["/bin/maze-api"]
