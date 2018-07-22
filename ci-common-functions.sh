
#
# Common CI/CD functions library.
#

# The following two functions are used to protect sensitive data from being
# printed to the console when shell execution tracing (-x) is turned on.
# Most of the CI scripts will "set -x" if "$ENABLE_DEBUG_BUILD_ENV" is set.

local __current_shell_flags=""

function disable_execution_tracing() {
   __current_shell_flags="$-"
   set +x
}

function reenable_execution_tracting() {
   [[ "$__current_shell_flags" =~ "x" ]] && set -x
}

