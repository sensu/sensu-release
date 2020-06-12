FROM alpine:3.12.0
MAINTAINER Sensu Engineering <engineering@sensu.io>

ENV ARTIFACTS_PATH="/artifacts"
ENV TOOLS_PATH="/tools"

RUN mkdir -p $ARTIFACTS_PATH
RUN mkdir -p $TOOLS_PATH

# Copy sensu-release tools to /tools
COPY ci-common-functions.sh $TOOLS_PATH
COPY post-packages-s3.sh $TOOLS_PATH
COPY go.mod $TOOLS_PATH
COPY go.sum $TOOLS_PATH
ADD cmd $TOOLS_PATH/cmd

ENV BUILD_TIME_PKGS="go"
ENV RUN_TIME_PKGS="openssh-client git aws-cli openssl curl bash sudo"

# Install packages
RUN apk add --no-cache $BUILD_TIME_PKGS $RUN_TIME_PKGS

# Build tools
RUN cd $TOOLS_PATH && go build ./cmd/circleci-logs

# Cleanup
RUN apk del $BUILD_TIME_PKGS
RUN rm -rf $TOOLS_PATH/go.mod $TOOLS_PATH/go.sum $TOOLS_PATH/cmd

# Add circleci user
RUN adduser -u 3434 -D -s /bin/bash circleci && \
    echo 'circleci ALL=NOPASSWD: ALL' >> /etc/sudoers.d/50-circleci && \
    echo 'Defaults    env_keep += "DEBIAN_FRONTEND"' >> /etc/sudoers.d/env_keep && \
    sudo -u circleci mkdir /home/circleci/project

# Fix ownership
RUN chown -R circleci:circleci $ARTIFACTS_PATH
RUN chown -R circleci:circleci $TOOLS_PATH

USER circleci
ENV HOME /home/circleci
ENV BASH_ENV $HOME/.bash_env
RUN touch $BASH_ENV
WORKDIR $HOME/project

CMD /bin/bash
