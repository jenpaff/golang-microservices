# The `./do` script

The idea of the `./do` script is to provide all development related tasks in a runnable shell script. Those tasks are typically:

* building the binary
* running tests
* linting the code
* reading credentials from a vault

Another benefit of `./do` scripts is that you can run the very same code you use locally in your pipeline. That way you can test your pipeline steps locally before running them in a pipeline. Also writing shell scripts in a syntax highlighted shell script is more convenient than having them inlined in a `.yaml` file.

> Although calling them `./do` scripts, we recommend adding a suffix `.sh` to the script to enable syntax highlighting in your editors.

## Dispatching Tasks

Our `./do.sh` script example provides a convenient way of adding tasks. All you have to do to add a new task is adding the following snippet:

```bash
## some-task-name: description
function task_some_task_name {
  param1=$1
  param2=$2
  param3=$3
  
  # ... your task ...
}
```

The dispatcher in the `./do` script will automatically dispatch the parameter given to the script as a first parameter and invoke the function with the remaining arguments:

```shell
./do.sh some-task-name param1 param2 param2
```

## Using `./do` Scripts in your GitHub Actions

Using the same script in your pipeline and on your workstation allows you for easy local testing before running tasks in your pipeline. E.g.

```yaml
on: [push]
name: Build
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Install Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.16.x
      - name: Checkout code
        uses: actions/checkout@v2
      - name: Build
        run: ./do build
```

We also think that this makes the pipeline easier to read. 

## TODOs

- [ ] add idea of self-documentation (man page)
- [ ] add auto-complete snippet for bash / zsh
