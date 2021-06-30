#!/usr/bin/env bash

set -euo pipefail

red='\033[1;31m'
green='\033[1;32m'
normal='\033[0m'

binary_path='webservice'
image_name='golang-microservices'
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

## lint : will run golangci linter
function task_lint() {
    if ! [ -x "$(command -v golangci-lint)" ]; then
        echo "Fetching linter..."
        pushd /tmp > /dev/null
        go get github.com/golangci/golangci-lint/cmd/golangci-lint
        popd > /dev/null
    fi

    golangci-lint run
}

## test : run all tests
function task_test {
    task_lint
    task_generate_all_mocks

    tags=${1:-"unit"}

    assert_ginkgo

    echo "Starting tests..."

    CONFIG_PATH="$(pwd)/config" ginkgo -r -tags="$tags" --randomizeAllSpecs --randomizeSuites --trace --progress -keepGoing --cover ./...
}

## integration-test: will spin up a sb and integration tests server and run integration tests
function task_integration_test {
  task_run_db

  assert_ginkgo

  CONFIG_PATH="$(pwd)/config" ginkgo -r -tags=integration --randomizeAllSpecs --randomizeSuites --trace --progress -keepGoing ./...
}

## container-test: will spin up local containers and run container tests
function task_container_test {
  task_build_container
  task_run_db

  CONFIG_PATH="$(pwd)/config" ginkgo -r -tags=container --randomizeAllSpecs --randomizeSuites --trace --progress -keepGoing ./...
}

## test-coverage : generate overall test coverage report and show in browser
function task_test_coverage {
  # our task_test creates a .coverprofile in each package
  # this task concatenates those profiles and generates
  # a overall .coverprofile.
  # finally, it outputs a report on the commandline and in the browser (if run in terminal)

  rm -f coverage/coverage.txt
  mkdir -p coverage
  touch coverage/coverage.txt
  echo "mode: set" >> coverage/coverage.txt

  task_test

  for d in $(find . | grep '.coverprofile'); do
    tail -n +2 ${d} >> coverage/coverage.txt
  done

  go tool cover -func=coverage/coverage.txt

  # only show coverage in browser, if run within a terminal
  if [ -t 1 ] ; then
    # see https://blog.golang.org/cover
    go tool cover -html=coverage/coverage.txt
  fi
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
  docker rm -f go-postgres
  task_build_container

  docker-compose up -d
  sleep 20
  task_build_migrations
  ./migration "localhost" "5432" "golangservice" "false" "postgres" "password"
}

## run-db : start local postgres database
function task_run_db {
    docker network create go-service 2>/dev/null || true
    green "Removing old database"
    docker rm -f go-postgres
    green "Pulling image"
    docker pull postgres:11

    docker run -d \
        -e POSTGRES_DB="golangservice" \
        -e POSTGRES_USER="postgres" \
        -e POSTGRES_PASSWORD="password" \
        -e POSTGRES_HOST_AUTH_METHOD="trust" \
        --name="go-postgres" \
        --network="go-service" \
        -p 5432:5432 \
        -m 128m \
        postgres:11
    green "Wait 2 seconds until database is up"
    sleep 2
    task_build_migrations
    green "Running migrations"
    ./migration "localhost" "5432" "golangservice" "false" "postgres" "password"
}

## build-migrations: build the go binary for our webservice
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
        pushd /tmp > /dev/null
        go get github.com/volatiletech/sqlboiler/v4
        go get github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql
        popd > /dev/null
    fi
    sqlboiler psql
}

## generate-all-mocks: generates all mocks for interfaces with mockgen annotation
function task_generate_all_mocks {
  echo "Generating all mocks..."
    if ! [ -x "$(command -v mockgen)" ]; then
        echo "Fetching mockgen..."
        pushd /tmp > /dev/null
        go get github.com/golang/mock/mockgen
        popd > /dev/null
    fi

    go generate ./...
}

# print out an auto-generated "man page" for this script
function task_usage {
    echo "Usage: $0"
    sed -n 's/^##//p' <"$0" | column -t -s ':' |  sed -E $'s/^/\t/' | sort
}

function green {
  echo -e "${green}$1${normal}"
}

function assert_ginkgo {
    if ! [ -x "$(command -v ginkgo)" ]; then
        echo "Installing ginkgo cli..."
        pushd /tmp > /dev/null
        GO111MODULE=on go get github.com/onsi/ginkgo/ginkgo
        popd > /dev/null
    fi
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