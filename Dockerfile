FROM golang:1.18.7-buster AS builder

ARG VERSION=dev

WORKDIR /go/src/app
COPY go.* ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN mkdir bin/ && go build -o bin/ -ldflags=-X=main.version=${VERSION} ./cmd/...

FROM debian:buster-slim

COPY --from=builder /go/src/app/bin/cmd /go/bin/api

EXPOSE 80/tcp
ENTRYPOINT ["/go/bin/api"]