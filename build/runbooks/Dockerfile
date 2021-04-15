FROM golang:buster as builder

RUN mkdir /build
ADD cmd/runbooks /build/cmd/runbooks
ADD internal /build/internal
ADD go.mod /build/go.mod
ADD go.sum /build/go.sum

RUN cd /build/cmd/runbooks && go build

FROM debian:buster
COPY --from=builder /build/cmd/runbooks/runbooks /usr/bin/runbooks
RUN chmod +x /usr/bin/runbooks
ENTRYPOINT ["/usr/bin/runbooks","--config","/var/conf/config.yml"]