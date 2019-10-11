#!/bin/bash

# $ publish-packages.sh sha branch semver
# $ publish-packages.sh 39d2bf22c416f9a78f22c4817dc74ab62d7f4d63 release_5.11 5.11.1
git_branch_no_slashes="$(echo "$2" | sed -e 's:/:_:g')"
backend_postfix="$(aws s3 ls s3://sensu-ci-builds/$git_branch_no_slashes/ | grep $1 | awk '{print $2}')"
aws s3 sync s3://sensu-ci-builds/$git_branch_no_slashes/${backend_postfix}deb s3://sensu.io/sensu-go/$3 --acl public-read
aws s3 sync s3://sensu-ci-builds/$git_branch_no_slashes/${backend_postfix}msi s3://sensu.io/sensu-go/$3 --acl public-read
aws s3 sync s3://sensu-ci-builds/$git_branch_no_slashes/${backend_postfix}rpm s3://sensu.io/sensu-go/$3 --acl public-read
aws s3 cp s3://sensu-ci-builds/$git_branch_no_slashes/${backend_postfix}goreleaser s3://sensu.io/sensu-go/$3 --recursive --exclude "*" --include "*.txt" --include "*.tar.gz" --include "*.zip" --acl public-read
aws s3 cp s3://sensu-ci-builds/$git_branch_no_slashes/${backend_postfix}log.txt s3://sensu.io/sensu-go/$3/log.txt