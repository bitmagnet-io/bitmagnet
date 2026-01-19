FROM golang:1.25.1-alpine3.22 AS base

RUN apk --update add \
  gcc \
  musl-dev \
  git \
  && rm -rf /var/cache/apk/*

FROM base AS dev

RUN apk --update add \
  chromium \
  curl \
  go-task \
  golangci-lint \
  iproute2-ss \
  jekyll \
  nodejs \
  npm \
  protobuf \
  ruby-bundler \
  starship \
  sudo \
  ruby \
  vim \
  make \
  zsh \
  zsh-autosuggestions \
  zsh-history-substring-search \
  zsh-syntax-highlighting \
  && rm -rf /var/cache/apk/*

RUN ln -s /usr/bin/go-task /usr/bin/task

RUN npm install -g prettier

# Create a new user and group
ARG USERNAME=bitmagnet
ARG UID=1000
ARG GID=1000

RUN addgroup -g "${GID}" "${USERNAME}" && \
  adduser -D -u "${UID}" -G "${USERNAME}" -s /bin/sh "${USERNAME}" && \
  mkdir -p /etc/sudoers.d && \
  echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/"${USERNAME}" && \
  chmod 0440 /etc/sudoers.d/"${USERNAME}"

USER ${USERNAME}

RUN touch /home/${USERNAME}/.bashrc

RUN mkdir -p /home/${USERNAME}/.local/bin

RUN mkdir -p /home/${USERNAME}/code

RUN chown -R ${USERNAME}:${USERNAME} /home/${USERNAME}

RUN mkdir -p /home/${USERNAME}/go/pkg/mod/cache && \
  mkdir -p /home/${USERNAME}/.cache && \
  mkdir -p /home/${USERNAME}/gotools

COPY go.mod go.sum /home/${USERNAME}/gotools/

ENV GOMODCACHE=/home/${USERNAME}/go/pkg/mod/cache

RUN cd /home/${USERNAME}/gotools && go mod download \
  && go install mvdan.cc/gofumpt@$(go list -m -f '{{.Version}}' mvdan.cc/gofumpt) \
  && go install google.golang.org/protobuf/cmd/protoc-gen-go@$(go list -m -f '{{.Version}}' google.golang.org/protobuf) \
  && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(go list -m -f '{{.Version}}' google.golang.org/grpc/cmd/protoc-gen-go-grpc) \
  && go install golang.org/x/tools/gopls@$(go list -m -f '{{.Version}}' golang.org/x/tools/gopls) \
  && go install github.com/go-delve/delve/cmd/dlv@$(go list -m -f '{{.Version}}' github.com/go-delve/delve) \
  && rm -rf /home/${USERNAME}/gotools

ENTRYPOINT ["tail", "-f", "/dev/null"]

FROM base AS build

RUN mkdir /build

COPY . /build

WORKDIR /build

RUN go build -ldflags "-s -w -X github.com/bitmagnet-io/bitmagnet/internal/version.GitTag=$(git describe --tags --always --dirty)"

FROM alpine:3.20 AS dist

COPY --from=build /build/bitmagnet /usr/bin/bitmagnet

ENTRYPOINT ["bitmagnet"]
