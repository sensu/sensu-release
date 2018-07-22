#!/usr/bin/env bash

set -e

readonly SCRIPT_DIR="$(dirname "$0")"
. $SCRIPT_DIR/ci-common-functions.sh

test -n "$DEBUG_BUILD_ENV" && set -x

readonly OPENSSL="${OPENSSL:-/usr/bin/openssl}"

function usage() {
   echo "Usage: $0 [ encrypted_variables_file ] [ password_variable ] [ decrypted_variables_file ] " >&2
}

if [[ -z "$1" ]]; then
   usage
   exit 1
fi
readonly encrypted_var_file="$1"

if [[ -z "$2" ]]; then
   usage
   exit 1
fi
readonly decryption_password_variable="$2"

if [[ -z "$3" ]]; then
   usage
   exit 1
fi
readonly decrypted_vars_file="$3"

disable_execution_tracing
decryption_password="${!decryption_password_variable}"
$OPENSSL aes-256-cbc -d -in "$encrypted_var_file" -out "$decrypted_vars_file" -k "$decryption_password"
reenable_execution_tracing

exit 0
