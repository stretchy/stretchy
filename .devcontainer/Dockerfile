FROM golang:1.15

WORKDIR /srv

COPY --from=golangci/golangci-lint:v1.31 /usr/bin/golangci-lint /usr/bin/golangci-lint
COPY --from=docker:stable /usr/local/bin/docker /usr/local/bin/docker
