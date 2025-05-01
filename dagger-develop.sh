#!/bin/bash

# This script runs 'dagger develop' in all modules in the repository, updating
# all dependencies and the dagger engine version. It must be ran in the repository
# root directory.

# Notes:
# dagger develop --recursive may be useful here, but it seems to require a
# top-level module with all of our modules as children.

latest_version=$(curl https://api.github.com/repos/dagger/dagger/releases -s | jq -r .[].tag_name | grep '^v[0-9]\.[0-9]*\.[0-9]*$' | head -n 1)

current_version=$(dagger version | cut -d ' ' -f 2)

if [ "$latest_version" != "$current_version" ]; do
    echo "Dagger engine is not the latest."
    echo "Current Version: $current_version"
    echo "Latest Version: $latest_version"
    exit 1
fi

for modulePath in $(pwd)/*/; do
    # skip dirs that aren't dagger modules
    if [ -f "$modulePath/dagger.json" ]; then
        dagger develop --sdk=go --mod "$modulePath"
    fi
done