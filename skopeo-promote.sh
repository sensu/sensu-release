#!/bin/sh

# This script will generate auth.json and store it in $SKOPEO_CONFIG_DIR for use
# with other Skopeo commands.
#
# If you already have used Skopeo locally you do not need to use this script.
# Instead, set SKOPEO_CONFIG_DIR to the directory on your local system which
# contains auth.json (likely ~/.config/containers).

usage() {
    echo "usage:"
    echo "  skopeo-copy.sh ci-to-testing <major_version> <minor_version> <patch_version> <git_sha>"
    echo "  skopeo-copy.sh ci-to-stable <major_version> <minor_version> <patch_version> <git_sha>"
    echo "  skopeo-copy.sh testing-to-stable <major_version> <minor_version> <patch_version>"
    exit 1
}

error() {
    echo "error: ${1}"
    echo
    usage
}

# SKOPEO_CONFIG_DIR is the path to the local config directory for Skopeo
SKOPEO_CONFIG_DIR="${PWD}/.config/containers"

SENSU_CI_DOCKER_REPO="docker://docker.io/sensu/sensu-ci"
SENSU_DOCKER_REPO="docker://docker.io/sensu/sensu"
SENSU_RHEL_DOCKER_REPO="docker://docker.io/sensu/sensu-rhel"
SENSU_TESTING_DOCKER_REPO="docker://docker.io/sensu/sensu-testing"

COMMAND="$1"
if [ "${COMMAND}" = "" ]; then
    error "missing command"
fi

MAJOR_VERSION="$2"
if [ "${MAJOR_VERSION}" = "" ]; then
    error "missing major version"
fi

MINOR_VERSION="$3"
if [ "${MINOR_VERSION}" = "" ]; then
    error "missing minor version"
fi

PATCH_VERSION="$4"
if [ "${PATCH_VERSION}" = "" ]; then
    error "missing patch version"
fi

GIT_SHA="$5"
case $COMMAND in
    ci-to-stable|ci-to-testing)
        if [ "$GIT_SHA" = "" ]; then
            error "missing git sha"
        fi
        ;;
esac

VERSION="${MAJOR_VERSION}.${MINOR_VERSION}.${PATCH_VERSION}"
SHORT_VERSION="${MAJOR_VERSION}.${MINOR_VERSION}"
ALPINE_TAG="alpine"
RHEL7_TAG="rhel7"

skopeo_copy() {
    src="$1"
    dst="$2"

    echo "running: skopeo copy -a ${src} to ${dst}"
    docker run \
        -v "${SKOPEO_CONFIG_DIR}:/tmp/config" \
        -e REGISTRY_AUTH_FILE="/tmp/config/auth.json" \
        -it quay.io/skopeo/stable:latest \
        copy -a "${src}" "${dst}"
}

ci_to_stable() {
    # create alpine tags
    skopeo_copy "${SENSU_CI_DOCKER_REPO}:${GIT_SHA}-${ALPINE_TAG}" "${SENSU_DOCKER_REPO}:${ALPINE_TAG}"
    skopeo_copy "${SENSU_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_DOCKER_REPO}:${VERSION}"
    skopeo_copy "${SENSU_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_DOCKER_REPO}:${SHORT_VERSION}"
    skopeo_copy "${SENSU_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_DOCKER_REPO}:${MAJOR_VERSION}"
    skopeo_copy "${SENSU_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_DOCKER_REPO}:latest"

    # create rhel7 tags
    skopeo_copy "${SENSU_CI_DOCKER_REPO}:${GIT_SHA}-${RHEL7_TAG}" "${SENSU_RHEL_DOCKER_REPO}:${RHEL7_TAG}"
    skopeo_copy "${SENSU_RHEL_DOCKER_REPO}:${RHEL7_TAG}" "${SENSU_RHEL_DOCKER_REPO}:${VERSION}"
    skopeo_copy "${SENSU_RHEL_DOCKER_REPO}:${RHEL7_TAG}" "${SENSU_RHEL_DOCKER_REPO}:${SHORT_VERSION}"
    skopeo_copy "${SENSU_RHEL_DOCKER_REPO}:${RHEL7_TAG}" "${SENSU_RHEL_DOCKER_REPO}:${MAJOR_VERSION}"
    skopeo_copy "${SENSU_RHEL_DOCKER_REPO}:${RHEL7_TAG}" "${SENSU_RHEL_DOCKER_REPO}:latest"
}

ci_to_testing() {
    # create alpine tags
    skopeo_copy "${SENSU_CI_DOCKER_REPO}:${GIT_SHA}-${ALPINE_TAG}" "${SENSU_TESTING_DOCKER_REPO}:${ALPINE_TAG}"
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_TESTING_DOCKER_REPO}:${VERSION}-${ALPINE_TAG}"
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_TESTING_DOCKER_REPO}:${SHORT_VERSION}-${ALPINE_TAG}"
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_TESTING_DOCKER_REPO}:${MAJOR_VERSION}-${ALPINE_TAG}"

    # create rhel7 tags
    skopeo_copy "${SENSU_CI_DOCKER_REPO}:${GIT_SHA}-${RHEL7_TAG}" "${SENSU_TESTING_DOCKER_REPO}:${RHEL7_TAG}"
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${RHEL7_TAG}" "${SENSU_TESTING_DOCKER_REPO}:${VERSION}-${RHEL7_TAG}"
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${RHEL7_TAG}" "${SENSU_TESTING_DOCKER_REPO}:${SHORT_VERSION}-${RHEL7_TAG}"
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${RHEL7_TAG}" "${SENSU_TESTING_DOCKER_REPO}:${MAJOR_VERSION}-${RHEL7_TAG}"

    # default tags (uses alpine)
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_TESTING_DOCKER_REPO}:${VERSION}"
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_TESTING_DOCKER_REPO}:${SHORT_VERSION}"
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_TESTING_DOCKER_REPO}:${MAJOR_VERSION}"
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_TESTING_DOCKER_REPO}:latest"
}

testing_to_stable() {
    # create alpine tags
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${VERSION}-${ALPINE_TAG}" "${SENSU_DOCKER_REPO}:${ALPINE_TAG}"
    skopeo_copy "${SENSU_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_DOCKER_REPO}:${VERSION}"
    skopeo_copy "${SENSU_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_DOCKER_REPO}:${SHORT_VERSION}"
    skopeo_copy "${SENSU_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_DOCKER_REPO}:${MAJOR_VERSION}"
    skopeo_copy "${SENSU_DOCKER_REPO}:${ALPINE_TAG}" "${SENSU_DOCKER_REPO}:latest"

    # create rhel7 tags
    skopeo_copy "${SENSU_TESTING_DOCKER_REPO}:${VERSION}-${RHEL7_TAG}" "${SENSU_RHEL_DOCKER_REPO}:${RHEL7_TAG}"
    skopeo_copy "${SENSU_RHEL_DOCKER_REPO}:${RHEL7_TAG}" "${SENSU_RHEL_DOCKER_REPO}:${VERSION}"
    skopeo_copy "${SENSU_RHEL_DOCKER_REPO}:${RHEL7_TAG}" "${SENSU_RHEL_DOCKER_REPO}:${SHORT_VERSION}"
    skopeo_copy "${SENSU_RHEL_DOCKER_REPO}:${RHEL7_TAG}" "${SENSU_RHEL_DOCKER_REPO}:${MAJOR_VERSION}"
    skopeo_copy "${SENSU_RHEL_DOCKER_REPO}:${RHEL7_TAG}" "${SENSU_RHEL_DOCKER_REPO}:latest"
}

case $COMMAND in
    ci-to-stable)
        ci_to_stable
        ;;
    ci-to-testing)
        ci_to_testing
        ;;
    testing-to-stable)
        testing_to_stable
        ;;
    *)
        error "invalid command"
        ;;
esac
