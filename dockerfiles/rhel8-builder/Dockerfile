ARG GORELEASER_VERSION=0.141.0
ARG CIRCLECI_LOGS_VERSION=69295d909442a39dac5b8ec386b71dba607cef3d

FROM registry.access.redhat.com/ubi8/python-38:1 as python

# Install aws-cli
RUN pip3 install https://github.com/boto/botocore/archive/v2.tar.gz
RUN pip3 install https://github.com/aws/aws-cli/archive/2.0.42.tar.gz

FROM registry.access.redhat.com/ubi8/go-toolset:1.13.4

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG GORELEASER_VERSION
ARG CIRCLECI_LOGS_VERSION

LABEL name="sensu/sensu-release" \
      maintainer="engineering@sensu.io"

RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM"

USER root

# Install python3.8
COPY --from=python /opt/app-root/bin/python3.8 /opt/app-root/bin/python3.8
COPY --from=python /opt/app-root/bin/python3.8 /opt/app-root/bin/python3
COPY --from=python /opt/app-root/bin/python3.8 /opt/app-root/bin/python
COPY --from=python /opt/app-root/lib/python3.8 /opt/app-root/lib/python3
COPY --from=python /opt/app-root/lib64/python3.8 /opt/app-root/lib64/python3
COPY --from=python /opt/app-root/bin/pip /opt/app-root/bin/pip
COPY --from=python /opt/app-root/bin/pip /opt/app-root/bin/pip3
COPY --from=python /opt/app-root/bin/pip /opt/app-root/bin/pip3.8

# Install packages
ENV SYSTEM_PKGS="openssl sudo"

RUN INSTALL_PKGS="${SYSTEM_PKGS}" && \
  yum install -y --setopt=tsflags=nodocs $INSTALL_PKGS && \
  rpm -V $INSTALL_PKGS && \
  yum -y clean all --enablerepo='*'

# Set envvars & make dirs
ENV WORKSPACE_PATH="/workspace"
ENV ARTIFACTS_PATH="${WORKSPACE_PATH}/artifacts"
ENV LOGS_PATH="${ARTIFACTS_PATH}/logs"
ENV TOOLS_PATH="/tools"

RUN mkdir -p $WORKSPACE_PATH
RUN mkdir -p $ARTIFACTS_PATH
RUN mkdir -p $LOGS_PATH
RUN mkdir -p $TOOLS_PATH

# Copy sensu-release tools to /tools
COPY ci-common-functions.sh $TOOLS_PATH
COPY post-packages-s3.sh $TOOLS_PATH
COPY go.mod $TOOLS_PATH
COPY go.sum $TOOLS_PATH
COPY docker2golang.sh $TOOLS_PATH
ADD cmd $TOOLS_PATH/cmd

# Install goreleaser
RUN export PLATFORM="$($TOOLS_PATH/docker2golang.sh $TARGETPLATFORM)" && \
  curl -sSL -o goreleaser.tar.gz https://s3-us-west-2.amazonaws.com/sensu.io/goreleaser-${GORELEASER_VERSION}.${PLATFORM}.tar.gz
RUN tar -C /usr/local/bin -zxf goreleaser.tar.gz
RUN rm goreleaser.tar.gz

# Install circleci-logs
RUN export PLATFORM="$($TOOLS_PATH/docker2golang.sh $TARGETPLATFORM)" && \
  curl -sSL -o circleci-logs.tar.gz https://s3-us-west-2.amazonaws.com/sensu.io/circleci-logs-${CIRCLECI_LOGS_VERSION}.${PLATFORM}.tar.gz
RUN tar -C /usr/local/bin -zxf circleci-logs.tar.gz
RUN rm circleci-logs.tar.gz

# Add circleci user
RUN useradd --uid=3434 --user-group --create-home circleci
RUN echo 'circleci ALL=NOPASSWD: ALL' >> /etc/sudoers.d/50-circleci
RUN echo 'Defaults    env_keep += "DEBIAN_FRONTEND"' >> /etc/sudoers.d/env_keep
RUN sudo -u circleci mkdir /home/circleci/project

# Fix ownership
RUN chown -R circleci:circleci $WORKSPACE_PATH
RUN chown -R circleci:circleci $ARTIFACTS_PATH
RUN chown -R circleci:circleci $LOGS_PATH
RUN chown -R circleci:circleci $TOOLS_PATH

# Fix BASH_ENV
RUN sed -e '/unset BASH_ENV PROMPT_COMMAND ENV/ s/^#*/#/' -i $APP_ROOT/etc/scl_enable

USER circleci
ENV HOME /home/circleci
ENV TERM xterm-color
WORKDIR $HOME/project
