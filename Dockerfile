FROM alpine:latest AS certificates

RUN apk --update add ca-certificates

FROM scratch

COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY stretchy /

ENTRYPOINT ["/stretchy"]