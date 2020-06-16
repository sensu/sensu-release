FROM registry.access.redhat.com/ubi8/go-toolset:1.13.4

LABEL name="sensu/sensu-release" \
      maintainer="engineering@sensu.io"

USER root

ENV WORKSPACE_PATH="/workspace"
ENV ARTIFACTS_PATH="${WORKSPACE_PATH}/artifacts"

RUN mkdir -p $WORKSPACE_PATH
RUN mkdir -p $ARTIFACTS_PATH

# Install packages
ENV SYSTEM_PKGS="openssl sudo"
ENV PYTHON_PKGS="python38"

RUN INSTALL_PKGS="${SYSTEM_PKGS} ${PYTHON_PKGS}" && \
  yum install -y --setopt=tsflags=nodocs $INSTALL_PKGS && \
  rpm -V $INSTALL_PKGS && \
  yum -y clean all --enablerepo='*'

ENV GO111MODULE="on"
ENV GOPROXY="https://proxy.golang.org"

# Install goreleaser
RUN rpm -Uvh "https://github.com/goreleaser/goreleaser/releases/download/v0.136.0/goreleaser_amd64.rpm"

# Add circleci user
RUN useradd --uid=3434 --user-group --create-home circleci
RUN echo 'circleci ALL=NOPASSWD: ALL' >> /etc/sudoers.d/50-circleci
RUN echo 'Defaults    env_keep += "DEBIAN_FRONTEND"' >> /etc/sudoers.d/env_keep
RUN sudo -u circleci mkdir /home/circleci/project

# Fix ownership
RUN chown -R circleci:circleci $WORKSPACE_PATH
RUN chown -R circleci:circleci $ARTIFACTS_PATH

# Fix BASH_ENV
RUN sed -e '/unset BASH_ENV PROMPT_COMMAND ENV/ s/^#*/#/' -i $APP_ROOT/etc/scl_enable

USER circleci
ENV HOME /home/circleci
WORKDIR $HOME/project