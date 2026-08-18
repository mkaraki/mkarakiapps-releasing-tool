#!/bin/sh
set -eu

cd /github/workspace

command="${INPUT_COMMAND:-}"

case "$command" in
  "pr-version-validator")
    echo "Running pr-version-validator..."
    exec /usr/bin/pr-version-validator \
      ${INPUT_TARGET_BRANCH:+--target-branch "$INPUT_TARGET_BRANCH"} \
      ${INPUT_VERSION_FILE:+--version-file "$INPUT_VERSION_FILE"} \
      ${INPUT_VERSION_FILE_FORMAT:+--version-file-format "$INPUT_VERSION_FILE_FORMAT"} \
      ${INPUT_LOCAL_PATCH_PREFIX:+--local-patch-prefix "$INPUT_LOCAL_PATCH_PREFIX"}
    ;;
  *)
    echo "Unknown command: $command"
    exit 1
    ;;
esac