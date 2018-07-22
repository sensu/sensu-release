#!/usr/bin/env bash

set -e
test -n "$DEBUG_BUILD_ENV" && set -x

readonly SCRIPT_DIR="$(dirname "$0")"
. $SCRIPT_DIR/ci-common-functions.sh

# The aws command is in awscli, because we install it ourselves.
readonly AWSCLI="${AWSCLI:-/usr/local/bin/aws}"
readonly GIT="${GIT:-/usr/bin/git}"
readonly SED="${SED:-/bin/sed}"

readonly AWS_S3_SENSU_CI_BUILDS_BUCKET="s3://sensu-ci-builds"

if [[ -z "$1" ]]; then
   echo "Usage: $0 [ deliverables directory ]" >&2
   exit 1
fi

readonly deliverables_dir="$1"

git_branch="$(cd $deliverables_dir && $GIT rev-parse --abbrev-ref HEAD)"
git_branch_no_slashes="$(echo "$git_branch" | $SED -e 's:/:_:g')"

build_date="$(cd $deliverables_dir && TZ='America/Los_Angeles' $GIT log -1 --format='%cd' --date=format:'%Y%m%d-%H%M' HEAD)"

disable_execution_tracing
export AWS_ACCESS_KEY_ID="$AWS_S3_SENSU_CI_BUILDS_ACCESS_KEY"
export AWS_SECRET_ACCESS_KEY="$AWS_S3_SENSU_CI_BUILDS_ACCESS_SECRET"
reenable_execution_tracing

$AWSCLI \
   s3 \
   cp \
   --recursive \
   --acl public-read \
   --metadata garbage-collect=1 \
   $deliverables_dir \
   $AWS_S3_SENSU_CI_BUILDS_BUCKET/$git_branch_no_slashes/$build_date

