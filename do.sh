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
  green "Generating models"
  task_generate_db_models
  green "Building go binary..."
  go get -d ./...
  go build -o ${binary_path} cmd/webservice/main.go
  green "Find binary in '${binary_path}'"
}

## build-container: builds a Docker image for our webservice
function task_build_container {
  green "Generating models"
  task_generate_db_models
  green "Start building Docker image..."
  docker image rm -f ${image_name}
  docker build --no-cache -t ${image_name} .
  green "Finished building Docker image '${image_name}'"
}

## run-container: run our webservice in a Docker container
function task_run_container {
  green "Removing old docker container"
  docker rm -f golang-microservices
  task_build_container
  green "Starting webservice in a Docker container"
  docker run -d --name ${container_name} -p 12345:12345 golang-microservices
  green "Container running - you can stop the container with 'docker stop ${container_name}'"
}

## run-db : start local postgres database
function task_run_db {
    green "Removing old database"
    docker rm -f local-db
    green "Pulling image"
    docker pull postgres:11

    docker run -d \
        -e POSTGRES_DB="postgres" \
        -e POSTGRES_USER="my-user" \
        -e POSTGRES_PASSWORD="my-users-password" \
        -e POSTGRES_HOST_AUTH_METHOD="trust" \
        --name="local-db" \
        -p 5432:5432 \
        -m 128m \
        postgres:11
    green "Wait 2 seconds until database is up"
    sleep 2
    task_build_migrations
    green "Running migrations"
    ./migration "localhost" "5432" "postgres" "false" "my-user" "my-users-password"
}

## build-migration: build the go binary for our webservice
function task_build_migrations {
  green "Building go binary for migrations ..."
  go get -d ./...
  go build -o migration cmd/migration/main.go
  green "Find binary in 'migration'"
}

## generate-db-models [override-files]: generates db-models for the tables in the database.
function task_generate_db_models {
    if ls ./persistence/models/* > /dev/null 2>/dev/null; then
      green "sqlboiler: Using cached DB models"
      return
    fi
    green "generating SQL model code"
    task_run_db
    if ! command -v sqlboiler > /dev/null ; then
        echo "Installing sqlboiler..."
        go get github.com/volatiletech/sqlboiler/v4
        go get github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql
    fi
    sqlboiler psql
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