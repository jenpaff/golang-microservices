#!/usr/bin/env bash

set -euo pipefail

red='\033[1;31m'
green='\033[1;32m'
normal='\033[0m'

binary_path='webservice'
image_name='golang-microservice'
container_name=${image_name}
local_port='12345'


## build: build the go binary for our webservice
function task_build {
  green "Building go binary..."
  go build -o ${binary_path} cmd/webservice/main.go
  green "Find binary in '${binary_path}'"
}

## build-container: builds a Docker image for our webservice
function task_build_container {
  green "Start building Docker image..."
  docker build --no-cache -t ${image_name} .
  green "Finished building Docker image '${image_name}'"
}

## run-container: run our webservice in a Docker container
function task_run_container {
  task_build_container
  green "Starting webservice in a Docker container"
  docker run -d --name ${container_name} -p 12345:12345 golang-microservices
  green "Container running - you can stop the container with 'docker stop ${container_name}'"
}

# print out an auto-generated "man page" for this script
function task_usage {
    echo "Usage: $0"
    sed -n 's/^##//p' <"$0" | column -t -s ':' |  sed -E $'s/^/\t/' | sort
}

function green {
  echo -e "${green}$1${normal}"
}

# read expected task as first CLI parameter
task=${1:-}

# remove first parameter from param list
shift || true

# translate a task with '-' into a function call with "_"
resolved_command=$(echo "task_${task}" | sed 's/-/_/g')

# check whether a function for the resolved command exists in this script
if [[ "$(LC_ALL=C type -t "${resolved_command}")" == "function" ]]; then
  # invoke function with remaining parameters
  pushd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null
  ${resolved_command} "$@"
else
  # if function could not be resolved, print out script usage
  task_usage
fi