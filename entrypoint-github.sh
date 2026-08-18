#!/bin/sh
set -eu

cd /github/workspace

command="${INPUT_COMMAND:-}"

case "$command" in
  "pr-version-validator")
    echo "Running pr-version-validator..."

    target_branch=$(printenv 'INPUT_TARGET-BRANCH' || true)
    version_file=$(printenv 'INPUT_VERSION-FILE' || true)
    version_file_format=$(printenv 'INPUT_VERSION-FILE-FORMAT' || true)
    local_patch_prefix=$(printenv 'INPUT_LOCAL-PATCH-PREFIX' || true)

    set -- /usr/bin/pr-version-validator

    if [ -n "$target_branch" ]; then
      set -- "$@" --target-branch "$target_branch"
    fi

    if [ -n "$version_file" ]; then
      set -- "$@" --version-file "$version_file"
    fi

    if [ -n "$version_file_format" ]; then
      set -- "$@" --version-file-format "$version_file_format"
    fi

    if [ -n "$local_patch_prefix" ]; then
      set -- "$@" --local-patch-prefix "$local_patch_prefix"
    fi

    exec "$@"
    ;;
  *)
    echo "Unknown command: $command"
    exit 1
    ;;
esac