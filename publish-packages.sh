#!/usr/bin/env bash

# $ publish-packages.sh sha branch semver
# $ publish-packages.sh 39d2bf22c416f9a78f22c4817dc74ab62d7f4d63 release/5.11 5.11.1
set -e
mkdir -p nupkgs
ci_bucket="sensu-ci-builds"
release_bucket="sensu.io"

# strip any underscores from the given git branch
git_branch_no_slashes="${2////_}"

# list directories under the directory that matches the branch name with no
# slashes and then find the directory that matches the given sha.
match="$(aws s3 ls s3://${ci_bucket}/${git_branch_no_slashes}/ | grep $1)"
dir_name_with_slash="$(echo $match | awk '{print $2}')"
dir_name="$(echo ${dir_name_with_slash%/})"


source_base_uri="s3://${ci_bucket}/${git_branch_no_slashes}/${dir_name}"
destination_base_uri="s3://${release_bucket}/sensu-go/${3}"

aws s3 sync "${source_base_uri}/build/deb" "${destination_base_uri}" --acl public-read
aws s3 sync "${source_base_uri}/build/msi" "${destination_base_uri}" --acl public-read
aws s3 sync "${source_base_uri}/build/rpm" "${destination_base_uri}" --acl public-read
aws s3 sync "${source_base_uri}/build/nupkg" "${destination_base_uri}" --acl public-read
aws s3 sync "${source_base_uri}/logs" "${destination_base_uri}" --acl public-read

aws s3 cp "${source_base_uri}/build/goreleaser" "${destination_base_uri}" --acl public-read --recursive --exclude "*" --include "*.txt" --include "*.tar.gz" --include "*.zip"

# Sync goreleaser dirs for non-standard builds
other_builds=("fips-openssl-1.0" "fips-openssl-1.1")
for build in ${other_builds[@]}; do
    aws s3 cp "${source_base_uri}/${build}/goreleaser" "${destination_base_uri}/${build}" --acl public-read --recursive --exclude "*" --include "*.txt" --include "*.tar.gz" --include "*.zip"
done

# Sync nupkg files to local nupkgs dir
aws s3 sync "${source_base_uri}/build/nupkg" nupkgs/
