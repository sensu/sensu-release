#!/usr/bin/env bash

set -e
test -n "$DEBUG_BUILD_ENV" && set -x

readonly CHMOD="${CHMOD:-/bin/chmod}"
readonly MKDIR="${MKDIR:-/bin/mkdir}"
readonly PYTHON="${PYTHON:-/usr/bin/python}"
readonly SUDO="${SUDO:-/usr/bin/sudo}"
readonly TAR="${TAR:-/bin/tar}"

readonly AWSCLI_INSTALL_DIR="/awscli-install"
readonly AWSCLI_VER="awscli-1.14.50"

readonly SETUP_CFG_EDIT="
[easy_install]
find_links = $SENSU_RELEASE_REPO/dist/${AWSCLI_VER}-deps
"

$MKDIR -p $AWSCLI_INSTALL_DIR
cd $AWSCLI_INSTALL_DIR
$TAR xfvz $SENSU_RELEASE_REPO/dist/${AWSCLI_VER}.tar.gz
cd $AWSCLI_VER

echo "$SETUP_CFG_EDIT" >> setup.cfg

$SUDO $PYTHON setup.py install

# For some reason, this egg gets installed with perms of 0640, which means mere
# mortals can't read it, which will make the awscli fail. So, fix that.
$SUDO $CHMOD 0644 /usr/local/lib/python2.7/dist-packages/pyasn1-0.4.3-py2.7.egg

exit 0
