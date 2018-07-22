
#
# Common CI/CD functions library.
#

# The following two functions are used to protect sensitive data from being
# printed to the console when shell execution tracing (-x) is turned on.
# Most of the CI scripts will "set -x" if "$ENABLE_DEBUG_BUILD_ENV" is set.

__original_shell_flags=""

function disable_execution_tracing() {
   __original_shell_flags="$-"
   set +x
}

function reenable_execution_tracing() {
   if [[ "$__original_shell_flags" =~ "x" ]]; then
      set -x
   fi
}

