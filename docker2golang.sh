#!/usr/bin/env bash
set -eo pipefail

case $1 in
    linux/386)
        echo "linux_386"
        ;;
    linux/amd64)
        echo "linux_amd64"
        ;;
    linux/arm64)
        echo "linux_arm64"
        ;;
    linux/arm/v6)
        echo "linux_arm_6"
        ;;
    linux/arm/v7)
        echo "linux_arm_7"
        ;;
    linux/ppc64le)
        echo "linux_ppc64le"
        ;;
    linux/s390x)
        echo "linux_s390x"
        ;;
    *)
        echo "no platform defined for $1"
        exit 1
        ;;
esac
