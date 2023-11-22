# The Debian version that the Ubuntu version is based on should be used as the
# value for DEBIAN_VERSION. E.g. Ubuntu 20.04 is based on Debian Bullseye.
ARG DEBIAN_VERSION="bullseye"
ARG UBUNTU_VERSION="20.04"
ARG GO_VERSION=1.21.3
ARG GOMPLATE_VERSION=3.9.0
ARG GORELEASER_VERSION=0.173.0
FROM golang:${GO_VERSION}-${DEBIAN_VERSION} as golang

# Build circleci-logs
COPY go.mod /tmp/sensu-release/
COPY go.sum /tmp/sensu-release/
ADD cmd /tmp/sensu-release/cmd
RUN cd /tmp/sensu-release && go build ./cmd/circleci-logs

FROM hairyhenderson/gomplate:v${GOMPLATE_VERSION} as gomplate
FROM goreleaser/goreleaser:v${GORELEASER_VERSION} as goreleaser
FROM ubuntu:${UBUNTU_VERSION} as toolchain

ENV PKGS="git wget xz-utils build-essential"

RUN apt update
RUN apt install -y $PKGS

ADD toolchain /toolchain

RUN make -C ./toolchain toolchain-linux/amd64 OUTPUT=/opt/toolchain -j3

FROM ubuntu:${UBUNTU_VERSION}

LABEL name="sensu/sensu-release" \
      maintainer="engineering@sensu.io"

ENV OBJ_PATH="/obj"
ENV WORKSPACE_PATH="/workspace"
ENV ARTIFACTS_PATH="${WORKSPACE_PATH}/artifacts"
ENV VAR_TMP_PATH="/var/tmp"
ENV RPM_LIB_PATH="/var/lib/rpm"
ENV LOGS_PATH="${ARTIFACTS_PATH}/logs"
ENV TOOLS_PATH="/tools"

RUN mkdir -p $OBJ_PATH
RUN mkdir -p $WORKSPACE_PATH
RUN mkdir -p $ARTIFACTS_PATH
RUN mkdir -p $VAR_TMP_PATH
RUN mkdir -p $RPM_LIB_PATH
RUN mkdir -p $LOGS_PATH
RUN mkdir -p $TOOLS_PATH

# Set timezone as some packages require it
ENV TZ=America/Vancouver \
    DEBIAN_FRONTEND=noninteractive

# Install packages
ENV PKGS="openssh-client git build-essential awscli bash sudo curl silversearcher-ag"

RUN apt update
RUN apt install -y $PKGS

# Install toolchain
COPY --from=toolchain /opt/toolchain /opt/toolchain

# Install go
COPY --from=golang /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:$PATH"
ENV GO111MODULE="on"
ENV GOPROXY="https://proxy.golang.org"

# Install publishing tools
COPY ci-common-functions.sh $TOOLS_PATH
COPY post-packages-s3.sh $TOOLS_PATH

# Install circleci-logs
COPY --from=golang /tmp/sensu-release/circleci-logs /usr/local/bin/

# Install goreleaser
COPY --from=goreleaser /usr/local/bin/goreleaser /usr/local/bin/goreleaser

# Install gomplate
COPY --from=gomplate /gomplate /usr/local/bin/gomplate

# Add circleci user
RUN adduser --uid 3434 --disabled-login --disabled-password --shell /bin/bash circleci
RUN echo 'circleci ALL=NOPASSWD: ALL' >> /etc/sudoers.d/50-circleci
RUN echo 'Defaults    env_keep += "DEBIAN_FRONTEND"' >> /etc/sudoers.d/env_keep

RUN sudo -u circleci mkdir /home/circleci/project

# Fix ownership
RUN chown -R circleci:circleci $WORKSPACE_PATH
RUN chown -R circleci:circleci $ARTIFACTS_PATH
RUN chown -R circleci:circleci $OBJ_PATH
RUN chown -R circleci:circleci $VAR_TMP_PATH
RUN chown -R circleci:circleci $RPM_LIB_PATH
RUN chown -R circleci:circleci $LOGS_PATH
RUN chown -R circleci:circleci $TOOLS_PATH
RUN chown -R circleci:circleci /tmp
RUN chown -R circleci:circleci /home/circleci

USER circleci
ENV HOME /home/circleci
WORKDIR $HOME/project
