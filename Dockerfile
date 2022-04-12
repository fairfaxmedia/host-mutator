FROM golang:1.18-alpine AS build

ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR /root/go/src

COPY go.mod go.sum .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build

FROM alpine:3.15

COPY --from=build /root/go/src/host-mutator /usr/local/bin/host-mutator

ENTRYPOINT ["/usr/local/bin/host-mutator"]
