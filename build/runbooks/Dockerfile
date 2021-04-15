FROM golang:buster as builder

WORKDIR /build
ADD cmd/runbooks cmd/runbooks
ADD internal internal
ADD go.mod go.mod
ADD go.sum go.sum

RUN cd cmd/runbooks && go build

FROM debian:buster
COPY --from=builder /build/cmd/runbooks/runbooks /usr/bin/runbooks
RUN chmod +x /usr/bin/runbooks
ENTRYPOINT ["/usr/bin/runbooks","--config","/var/conf/config.yml"]