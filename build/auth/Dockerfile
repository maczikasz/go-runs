FROM golang:buster as builder

RUN mkdir /build
ADD cmd/auth /build/cmd/auth
ADD internal /build/internal
ADD go.mod /build/go.mod
ADD go.sum /build/go.sum

RUN cd /build/cmd/auth && go build

FROM debian:buster
COPY --from=builder /build/cmd/auth/auth /usr/bin/auth
RUN chmod +x /usr/bin/auth
ENTRYPOINT ["/usr/bin/auth","--config","/var/conf/config.yml"]