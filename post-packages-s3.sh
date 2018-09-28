#!/usr/bin/env bash

set -e
test -n "$DEBUG_BUILD_ENV" && set -x

readonly SCRIPT_DIR="$(dirname "$0")"
. $SCRIPT_DIR/ci-common-functions.sh

# The aws command is in awscli, because we install it ourselves.
readonly AWK="${AWK:-/usr/bin/awk}"
readonly AWSCLI="${AWSCLI:-/usr/local/bin/aws}"
readonly GIT="${GIT:-/usr/bin/git}"
readonly SED="${SED:-/bin/sed}"

readonly AWS_S3_SENSU_CI_BUILDS_BUCKET="sensu-ci-builds"

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
   $deliverables_dir \
   s3://$AWS_S3_SENSU_CI_BUILDS_BUCKET/$git_branch_no_slashes/$build_date

uploaded_build_artifacts="$($AWSCLI s3 ls --recursive s3://$AWS_S3_SENSU_CI_BUILDS_BUCKET/$git_branch_no_slashes/$build_date | $AWK '{print $4}')"

for obj in $uploaded_build_artifacts; do
   echo "Tagging $obj for artifact garbage collection..."
   $AWSCLI s3api put-object-tagging --bucket $AWS_S3_SENSU_CI_BUILDS_BUCKET \
     --key "$obj" \
     --tagging '{
                  "TagSet": [
                     {
                       "Key": "garbage-collect",
                       "Value": "1"
                      }
                   ]
                }'

   sleep 1
done

