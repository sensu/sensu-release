FROM golang:1.13.12-alpine3.12 as golang
FROM alpine:3.12.0

LABEL name="sensu/sensu-release" \
      maintainer="engineering@sensu.io"

ENV OBJ_PATH="/obj"
ENV WORKSPACE_PATH="/workspace"
ENV ARTIFACTS_PATH="${WORKSPACE_PATH}/artifacts"
ENV VAR_TMP_PATH="/var/tmp"
ENV RPM_LIB_PATH="/var/lib/rpm"

RUN mkdir -p $OBJ_PATH
RUN mkdir -p $WORKSPACE_PATH
RUN mkdir -p $ARTIFACTS_PATH
RUN mkdir -p $VAR_TMP_PATH
RUN mkdir -p $RPM_LIB_PATH

# Install packages
ENV BUILD_TIME_PKGS="curl ruby-dev tar gcc g++ py3-setuptools"
ENV RUN_TIME_PKGS="cmake python3 ruby openssh-client git make rpm dpkg dpkg-dev aws-cli bash sudo"

RUN apk add --no-cache $BUILD_TIME_PKGS $RUN_TIME_PKGS

# Install go
COPY --from=golang /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:$PATH"
ENV GO111MODULE="on"
ENV GOPROXY="https://proxy.golang.org"

# Install goreleaser
RUN curl -Lo /tmp/goreleaser.tar.gz https://github.com/goreleaser/goreleaser/releases/download/v0.136.0/goreleaser_Linux_x86_64.tar.gz
RUN tar -C /tmp -zxf /tmp/goreleaser.tar.gz
RUN mv /tmp/goreleaser /usr/bin/

# Install the packagecloud cli tool
RUN gem install package_cloud -v 0.3.05

# Symlink python to python3
RUN ln -s /usr/bin/python3 /usr/bin/python

# Cleanup
RUN rm -rf /tmp/*
RUN apk del $BUILD_TIME_PKGS

# Add circleci user
RUN adduser -u 3434 -D -s /bin/bash circleci
RUN echo 'circleci ALL=NOPASSWD: ALL' >> /etc/sudoers.d/50-circleci
RUN echo 'Defaults    env_keep += "DEBIAN_FRONTEND"' >> /etc/sudoers.d/env_keep
RUN sudo -u circleci mkdir /home/circleci/project

# Fix ownership
RUN chown -R circleci:circleci $WORKSPACE_PATH
RUN chown -R circleci:circleci $ARTIFACTS_PATH
RUN chown -R circleci:circleci $OBJ_PATH
RUN chown -R circleci:circleci $VAR_TMP_PATH
RUN chown -R circleci:circleci $RPM_LIB_PATH

USER circleci
ENV HOME /home/circleci
WORKDIR $HOME/project
