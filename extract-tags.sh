#!/bin/bash

REF_NAME="$1"
if [ -z "$REF_NAME" ]; then
    echo 'The <ref_name> paramter is required.'
    exit 1
fi

SEMVER_REGEX='^v?([0-9]+)\.([0-9]+)\.([0-9]+)$'
if ! [[ "$REF_NAME" =~ $SEMVER_REGEX ]]; then
    echo 'Only <major>.<minor>.<patch> semver release versions are supported.'
  exit 1
fi


SEMVER_COMPONENTS_STR=$(echo "$REF_NAME" | sed -E "s/$SEMVER_REGEX/\1 \2 \3/")
read -a SEMVER_COMPONENTS <<< "$SEMVER_COMPONENTS_STR"

echo "VER_MAJOR=${SEMVER_COMPONENTS[0]}"
echo "VER_MAJOR_MINOR=${SEMVER_COMPONENTS[0]}.${SEMVER_COMPONENTS[1]}"
echo "VER_MAJOR_MINOR_PATCH=${SEMVER_COMPONENTS[0]}.${SEMVER_COMPONENTS[1]}.${SEMVER_COMPONENTS[2]}"
