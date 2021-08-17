#!/bin/bash

env=$1
args=("$@")
package_ids=("${args[@]:1}")
environments=("local" "remote")

rootFilePath=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  pwd -P
)

if [[ "${environments[@]}" =~ "$env" ]]; then
  for package_id in "${package_ids[@]}"; do
    package_path=""

    for path in /instant/*; do
      if [[ -f "${path}/instant.json" ]] && [[ "$(sed -n "s/\"id\"://p" "${path}/instant.json" | sed -e 's/^[[:space:]]*//' | tr -d '[,\'\"])" == "${package_id}" ]]; then
        package_path="${path}"
        break
      fi
    done

    if [[ -n $package_path ]]; then
      echo -e "\n\nRunning Tests for package with id - ${package_id}\n\n"
      env-cmd -f ".env.${env}" cucumber-js "${package_path}/features/" --publish-quiet
    else
      echo -e "\n\nPackage with id ${package_id} has not been instantiated!\n\n"
    fi
  done
else
  echo 'Deployment Environment not specified! Should be "local" or "remote".'
fi
