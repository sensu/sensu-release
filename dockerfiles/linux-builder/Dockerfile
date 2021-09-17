ARG ALPINE_VERSION=3.14.1
ARG GO_VERSION=1.17.1
ARG GOMPLATE_VERSION=3.9.0
ARG GORELEASER_VERSION=0.173.0
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as golang

# Build circleci-logs
COPY go.mod /tmp/sensu-release/
COPY go.sum /tmp/sensu-release/
ADD cmd /tmp/sensu-release/cmd
RUN cd /tmp/sensu-release && go build ./cmd/circleci-logs

FROM hairyhenderson/gomplate:v${GOMPLATE_VERSION}-alpine as gomplate
FROM goreleaser/goreleaser:v${GORELEASER_VERSION} as goreleaser
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

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

# Install packages
ENV BUILD_TIME_PKGS="tar"
ENV RUN_TIME_PKGS="openssh-client git make aws-cli bash sudo curl jq"

RUN apk add --no-cache $BUILD_TIME_PKGS $RUN_TIME_PKGS

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
COPY --from=goreleaser /bin/goreleaser /usr/local/bin/goreleaser

# Install gomplate
COPY --from=gomplate /bin/gomplate /usr/local/bin/gomplate

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
RUN chown -R circleci:circleci $LOGS_PATH
RUN chown -R circleci:circleci $TOOLS_PATH
RUN chown -R circleci:circleci /tmp

USER circleci
ENV HOME /home/circleci
WORKDIR $HOME/project
