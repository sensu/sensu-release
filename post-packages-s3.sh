#!/bin/bash

set -e
test -n "$DEBUG_BUILD_ENV" && set -x

readonly SCRIPT_DIR="$(dirname "$0")"
. $SCRIPT_DIR/ci-common-functions.sh

# The aws command is in awscli, because we install it ourselves.
readonly AWS_S3_SENSU_CI_BUILDS_BUCKET="sensu-ci-builds"

if [[ -z "$1" ]]; then
   echo "Usage: $0 [ deliverables directory ]" >&2
   exit 1
fi

readonly deliverables_dir="$1"

git_branch="$CIRCLE_BRANCH"
git_branch_no_slashes="$(echo "$git_branch" | sed -e 's:/:_:g')"
git_sha="$CIRCLE_SHA1"

build_date="$COMMIT_DATE"
bucket_dir="$build_date"
bucket_dir+="_$git_sha"

disable_execution_tracing
export AWS_ACCESS_KEY_ID="$AWS_S3_SENSU_CI_BUILDS_ACCESS_KEY"
export AWS_SECRET_ACCESS_KEY="$AWS_S3_SENSU_CI_BUILDS_ACCESS_SECRET"
reenable_execution_tracing

aws \
   s3 \
   cp \
   --recursive \
   --acl public-read \
   $deliverables_dir \
   s3://$AWS_S3_SENSU_CI_BUILDS_BUCKET/$git_branch_no_slashes/$bucket_dir

uploaded_build_artifacts=$(cd $deliverables_dir && find * -type f)

for file in $uploaded_build_artifacts; do
   obj="${git_branch_no_slashes}/${bucket_dir}/${file}"
   echo "Tagging $obj for artifact garbage collection..."
   aws s3api put-object-tagging --bucket $AWS_S3_SENSU_CI_BUILDS_BUCKET \
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
