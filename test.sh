#!/bin/bash

env=$1
args=("$@")
package_ids=("${args[@]:1}")
environments=("local" "remote")

rootFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

# Retrieve cucumber feature files for packages to be tested and place them in temporary test folder.
# The feature files are contained in the instant docker volume, created when the services are started.
mkdir "$rootFilePath"/test

docker create --name test-helper -v instant:/instant busybox
docker cp test-helper:/instant ./test/.
docker rm test-helper

# Remove node_modules from the test folder, if any.
# Cucumber is being run in the instant root folder and the node_modules used should be the those in the root folder
rm -rf ./test/instant/node_modules

if [[ "${environments[@]}" =~ "$env" ]]; then
    for package_id in "${package_ids[@]}"; do
        package_path=""

        for path in "$rootFilePath"/test/instant/*; do
            if [[ -f "${path}/instant.json" ]] && [[ "$(sed -n "s/\"id\"://p" "${path}/instant.json" | sed -e 's/^[[:space:]]*//' | tr -d '[,\'\"])" == "${package_id}" ]]; then
                package_path="${path}"
                break
            fi
        done

        if [[ -n $package_path ]]; then
            echo -e "\n\nRunning Tests for package with id - ${package_id}\n\n"
            env-cmd -f ".env.${env}" cucumber-js -f progress-bar "${package_path}/features/"
        else
            echo -e "\n\nPackage with id ${package_id} has not been instantiated!\n\n"
        fi
    done
else
    echo 'Deployment Environment not specified! Should be "local" or "remote".'
fi

# Clean up
rm -rf "$rootFilePath/test"
