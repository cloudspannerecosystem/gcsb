FROM golang:1.16 as builder

WORKDIR /app/

COPY . .

ENV GOBIN="/root/bin"
ENV CGO_ENABLED=0
ENV GOFLAGS=-mod=vendor
ENV PATH="/usr/local/go/bin/:/root/bin:${PATH}"

RUN make build

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/gcsb /gcsb

ENTRYPOINT [ "/gcsb" ]
