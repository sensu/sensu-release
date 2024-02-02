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
aws s3 sync "${source_base_uri}/build/rpm" "${destination_base_uri}" --acl public-read
aws s3 sync "${source_base_uri}/build/msi" "${destination_base_uri}" --acl public-read
aws s3 sync "${source_base_uri}/build/nupkg" "${destination_base_uri}" --acl public-read
aws s3 sync "${source_base_uri}/logs" "${destination_base_uri}/logs"

aws s3 cp "${source_base_uri}/build" "${destination_base_uri}" --acl public-read --recursive --exclude "*" --include "*.tar.gz" --include "*.zip"

# Sync build dirs for non-standard builds
other_builds=("cgo" "build-fips")
for build in ${other_builds[@]}; do
    aws s3 cp "${source_base_uri}/${build}" "${destination_base_uri}/${build}" --acl public-read --recursive --exclude "*" --include "*.txt" --include "*.tar.gz" --include "*.zip"
done

# download checksum files, concatenate them, and then upload the checksums file
checksums_dir=$(mktemp -d)
aws s3 cp "${source_base_uri}" "${checksums_dir}" --recursive --exclude "*" --include "build/*.txt" --include "cgo/*.txt" --include "build-fips/*.txt"
checksums_file="${checksums_dir}/sensu-go_${3}_checksums.txt"
for build in ${other_builds[@]}; do
    build_dir="${checksums_dir}/${build}"
    tmp_file=$(mktemp)
    find "${build_dir}" -type f -exec awk '{gsub(/sensu-go/,"'$build'/sensu-go")}1' {} > $tmp_file \;
    find "${build_dir}" -type f -exec mv $tmp_file {} \;
done
find "${checksums_dir}" -type f -exec cat {} \; | sort -k2 > $checksums_file
aws s3 cp "${checksums_file}" "${destination_base_uri}/${checksums_filename}" --acl public-read

# Sync nupkg files to local nupkgs dir
aws s3 sync "${source_base_uri}/build/nupkg" nupkgs/
