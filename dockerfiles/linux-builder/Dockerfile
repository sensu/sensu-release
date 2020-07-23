FROM golang:1.13.12-alpine3.12 as golang
FROM sensu/sensu-release@sha256:bcd0e081ae40111c305a8e401dbafac64abad40b52bbbfefe06d4f419ea8ccd9 as publishing-tools
FROM alpine:3.12.0

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
ENV BUILD_TIME_PKGS="tar gcc g++ py3-setuptools"
ENV RUN_TIME_PKGS="cmake python3 openssh-client git make rpm dpkg dpkg-dev aws-cli bash sudo curl"

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

# Install gomplate
RUN curl -Lo /usr/bin/gomplate https://github.com/hairyhenderson/gomplate/releases/download/v3.7.0/gomplate_linux-amd64-slim
RUN chmod +x /usr/bin/gomplate

# Install the packagecloud-go cli tool
RUN curl -Lo /tmp/packagecloud-go.tar.gz https://github.com/amdprophet/packagecloud-go/releases/download/0.1.2/packagecloud-go_0.1.2_linux_amd64.tar.gz
RUN tar -C /tmp -zxf /tmp/packagecloud-go.tar.gz
RUN mv /tmp/packagecloud /usr/bin/

# Symlink python to python3
RUN ln -s /usr/bin/python3 /usr/bin/python

# Install publishing tools
COPY --from=publishing-tools /usr/lib/python3.8 /usr/lib/python3.8
COPY --from=publishing-tools /usr/bin/python3.8 /usr/bin/python3.8
COPY --from=publishing-tools /usr/bin/python3 /usr/bin/python3
COPY --from=publishing-tools $TOOLS_PATH/ci-common-functions.sh $TOOLS_PATH/ci-common-functions.sh
COPY --from=publishing-tools $TOOLS_PATH/circleci-logs $TOOLS_PATH/circleci-logs
COPY --from=publishing-tools $TOOLS_PATH/post-packages-s3.sh $TOOLS_PATH/post-packages-s3.sh

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

USER circleci
ENV HOME /home/circleci
WORKDIR $HOME/project
