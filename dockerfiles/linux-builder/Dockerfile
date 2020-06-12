FROM golang:1.13.12-alpine3.12 as golang
FROM alpine:3.12.0
MAINTAINER Sensu Engineering <engineering@sensu.io>

WORKDIR /app

ENV OBJ_PATH="/obj"
ENV DIST_PATH="/dist"

RUN mkdir -p $OBJ_PATH
RUN mkdir -p $DIST_PATH

ENV BUILD_TIME_PKGS="curl ruby-dev tar gcc g++ py3-setuptools"
ENV RUN_TIME_PKGS="cmake python3 ruby openssh-client git make rpm dpkg dpkg-dev aws-cli"

# Install packages
RUN apk add --no-cache $BUILD_TIME_PKGS $RUN_TIME_PKGS

# Install go
COPY --from=golang /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:$PATH"

# Install goreleaser
RUN curl -Lo /tmp/goreleaser.tar.gz https://github.com/goreleaser/goreleaser/releases/download/v0.128.0/goreleaser_Linux_x86_64.tar.gz
RUN tar -C /tmp -zxf /tmp/goreleaser.tar.gz
RUN mv /tmp/goreleaser /usr/bin/

# Install the packagecloud cli tool
RUN gem install package_cloud -v 0.3.05

# Symlink python to python3
RUN ln -s /usr/bin/python3 /usr/bin/python

# Cleanup
RUN rm -rf /tmp/*
RUN apk del $BUILD_TIME_PKGS
