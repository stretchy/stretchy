FROM golang:1.14.2

WORKDIR /srv

COPY --from=golangci/golangci-lint:v1.26 /usr/bin/golangci-lint /usr/bin/golangci-lint
COPY --from=docker:stable /usr/local/bin/docker /usr/local/bin/docker
