FROM golang:1.13 AS build
WORKDIR /build
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o host-mutator

FROM alpine:3.10
COPY --from=build /build/host-mutator /bin/host-mutator
CMD ["/bin/host-mutator"]