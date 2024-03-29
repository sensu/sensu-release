FROM alpine:3.14.1

LABEL name="sensu/sensu-release" \
      maintainer="engineering@sensu.io"

ENV TOOLS_PATH="/tools"

RUN mkdir -p $TOOLS_PATH

# Copy sensu-release tools to /tools
COPY go.mod $TOOLS_PATH
COPY go.sum $TOOLS_PATH
ADD cmd $TOOLS_PATH/cmd

# Install packages
ENV BUILD_TIME_PKGS="go"
ENV RUN_TIME_PKGS="openssh-client git bash sudo"

RUN apk add --no-cache $BUILD_TIME_PKGS $RUN_TIME_PKGS

# Build tools
ENV GO111MODULE="on"
ENV GOPROXY="https://proxy.golang.org"
RUN cd $TOOLS_PATH && go build ./cmd/packagecloudpruner

# Cleanup
RUN apk del $BUILD_TIME_PKGS
RUN rm -rf $TOOLS_PATH/go.mod $TOOLS_PATH/go.sum $TOOLS_PATH/cmd

# Add circleci user
RUN adduser -u 3434 -D -s /bin/bash circleci
RUN echo 'circleci ALL=NOPASSWD: ALL' >> /etc/sudoers.d/50-circleci
RUN echo 'Defaults    env_keep += "DEBIAN_FRONTEND"' >> /etc/sudoers.d/env_keep
RUN sudo -u circleci mkdir /home/circleci/project

# Fix ownership
RUN chown -R circleci:circleci $TOOLS_PATH

USER circleci
ENV HOME /home/circleci
ENV BASH_ENV $HOME/.bash_env
RUN touch $BASH_ENV
WORKDIR $HOME/project

CMD /bin/bash
