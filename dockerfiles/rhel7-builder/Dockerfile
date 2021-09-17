FROM registry.access.redhat.com/ubi7/go-toolset:1.15.14

LABEL name="sensu/sensu-release" \
      maintainer="engineering@sensu.io"

USER root

# Install packages
ENV SYSTEM_PKGS="openssl sudo"

RUN INSTALL_PKGS="${SYSTEM_PKGS}" && \
  yum install -y yum-utils && \
  yum install -y --setopt=tsflags=nodocs $INSTALL_PKGS && \
  rpm -V $INSTALL_PKGS

# Uninstall git and unstall newer git
RUN yum -y --setopt=tsflags=nodocs install http://mirror.centos.org/centos/7/os/x86_64/Packages/pcre2-10.23-2.el7.x86_64.rpm
RUN yum -y --setopt=tsflags=nodocs install http://mirror.centos.org/centos/7/os/x86_64/Packages/emacs-filesystem-24.3-23.el7.noarch.rpm
RUN yum -y --setopt=tsflags=nodocs install http://mirror.centos.org/centos/7/os/x86_64/Packages/libsecret-0.18.6-1.el7.x86_64.rpm
RUN yum -y --setopt=tsflags=nodocs install https://packages.endpoint.com/rhel/7/os/x86_64/endpoint-repo-1.7-1.x86_64.rpm
RUN yum -y --setopt=tsflags=nodocs install git
RUN yum -y clean all --enablerepo='*'

# Install aws-cli
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install

# Install goreleaser
RUN rpm -Uvh "https://github.com/goreleaser/goreleaser/releases/download/v0.136.0/goreleaser_amd64.rpm"

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
ADD cmd $TOOLS_PATH/cmd

# Build tools
ENV GO111MODULE="on"
ENV GOPROXY="https://proxy.golang.org"
RUN export PATH=$PATH:$(dirname $(find /opt -name go -type f)) && \
  cd $TOOLS_PATH && go build ./cmd/circleci-logs
RUN cp $TOOLS_PATH/circleci-logs /usr/local/bin

# Add circleci user
RUN useradd --uid=3434 --user-group --create-home circleci
RUN echo 'circleci ALL=NOPASSWD: ALL' >> /etc/sudoers.d/50-circleci
RUN echo 'Defaults    env_keep += "DEBIAN_FRONTEND"' >> /etc/sudoers.d/env_keep
RUN sudo -u circleci mkdir /home/circleci/project

# Fix ownership & permissions
RUN chown -R circleci:circleci $WORKSPACE_PATH
RUN chown -R circleci:circleci $ARTIFACTS_PATH
RUN chown -R circleci:circleci $APP_ROOT
RUN chown -R circleci:circleci $LOGS_PATH
RUN chown -R circleci:circleci $TOOLS_PATH
RUN chmod +x $APP_ROOT/etc/scl_enable

# Fix BASH_ENV
RUN sed -e '/unset BASH_ENV PROMPT_COMMAND ENV/ s/^#*/#/' -i $APP_ROOT/etc/scl_enable

USER circleci
ENV HOME /home/circleci
ENV TERM xterm-color
WORKDIR $HOME/project
