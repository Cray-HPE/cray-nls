#!/bin/bash
#
# Usage: bump-container-version.sh (major|minor|patch)
#
# Bump the version of the cray-nls container image in all the appropriate places
# in this repository. This includes the following:
#
#   - The `.version` file
#   - The container image version under `cray-service.containers.cray-nls.image.tag`
#     in the Helm chart's `values.yaml`.
#
# Then it also bumps the version of the Cray Helm chart itself. This script
# currently assumes the only change to the Helm chart is the version of the
# container image it's using, so it just bumps the patch version of the Helm
# chart.
#
CHARTS_PATH="$(dirname $0)/../charts/v1.0/cray-nls"
VALUES_FILE="$CHARTS_PATH/values.yaml"
CHART_FILE="$CHARTS_PATH/Chart.yaml"
VERSION_PROPERTY=".cray-service.containers.cray-nls.image.tag"

function print_usage() {
    cat << EOF
Usage: bump-container-version.sh (major|minor|patch)
EOF
}

function err_exit() {
    echo "$@" 1>&2
    exit 1
}

function get_container_version() {
    local dot_version
    local chart_container_version

    dot_version="$(cat .version)"
    chart_container_version="$(yq "$VERSION_PROPERTY" "$VALUES_FILE")"
    
    if [[ "$dot_version" != "$chart_container_version" ]]; then
        err_exit "Container version in .version file ($dot_version) does not match" \
            "version in chart values.yaml ($chart_container_version)."
    fi

    echo "$dot_version"
}

function update_container_version() {
    local new_version="$1"

    if [ -z $1 ]; then
        err_exit "update_container_version requires the version as an argument"
    fi

    echo "$new_version" > .version
    yq -i "$VERSION_PROPERTY = \"$new_version\"" "$VALUES_FILE"
}

function get_chart_version() {
    # For some reason, this does not work
    # yq ".version" "$CHART_FILE"
    grep "^version: " $CHART_FILE | awk '{print $2}'
}

function update_chart_version() {
    local new_version="$1"

    if [ -z $1 ]; then
        err_exit "update_chart_version requires the version as an argument"
    fi

    # For some reason, this does not work
    # yq ".version = '$new_version'" "$CHART_FILE"

    sed -i "s/^version: .*/version: $new_version/" "$CHART_FILE"
}

function discard_whitespace_changes() {
    # The next three lines are to avoid whitespace changes to the VALUES_FILE
    # Adapted from SO: https://stackoverflow.com/a/45486981
    git diff -U0 -w --ignore-blank-lines --no-color | git apply --cached --ignore-whitespace --unidiff-zero -
    # Discard whitespace changes
    git checkout .
    # Unstage changes which were staged by git apply
    git reset --mixed
}

if [[ $# -ne 1 ]]; then
    print_usage
    err_exit "$(basename $0) requires an argument."
fi
component="$1"

container_version=$(get_container_version)
container_major=$(awk -F '.' '{print $1}' <<<$container_version)
container_minor=$(awk -F '.' '{print $2}' <<<$container_version)
container_patch=$(awk -F '.' '{print $3}' <<<$container_version)
chart_version=$(get_chart_version)
chart_major=$(awk -F '.' '{print $1}' <<<$chart_version)
chart_minor=$(awk -F '.' '{print $2}' <<<$chart_version)
chart_patch=$(awk -F '.' '{print $3}' <<<$chart_version)

if [[ $component == "patch" ]]; then
    container_patch="$(( container_patch + 1 ))"
elif [[ $component == "minor" ]]; then
    container_minor="$(( container_minor + 1 ))"
elif [[ $component == "major" ]]; then
    container_major="$(( container_major + 1 ))"
fi

new_container_version="${container_major}.${container_minor}.${container_patch}"
new_chart_version="${chart_major}.${chart_minor}.$(( chart_patch + 1 ))"

echo "Updating container version from $container_version to $new_container_version"
echo "Updating chart version from $chart_version to $new_chart_version"

update_container_version "$new_container_version"
update_chart_version "$new_chart_version"

discard_whitespace_changes


