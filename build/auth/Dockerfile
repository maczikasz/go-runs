FROM golang:buster as builder

WORKDIR /build
ADD cmd/auth cmd/auth
ADD internal internal
ADD go.mod go.mod
ADD go.sum go.sum

RUN cd cmd/auth && go build

FROM debian:buster
COPY --from=builder /build/cmd/auth/auth /usr/bin/auth
RUN chmod +x /usr/bin/auth
ENTRYPOINT ["/usr/bin/auth","--config","/var/conf/config.yml"]