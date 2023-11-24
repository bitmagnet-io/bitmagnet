FROM golang:alpine3.18 AS build

RUN apk --update add \
    gcc \
    musl-dev \
    git

RUN mkdir /build

COPY . /build

WORKDIR /build

RUN go build -ldflags "-X github.com/bitmagnet-io/bitmagnet/internal/version.GitTag=$(git describe --tags --always --dirty)"

FROM alpine:3.18

RUN apk --update add \
    curl \
    iproute2-ss \
    && rm -rf /var/cache/apk/*

COPY --from=build /build/bitmagnet /usr/bin/bitmagnet

ENTRYPOINT ["bitmagnet"]
