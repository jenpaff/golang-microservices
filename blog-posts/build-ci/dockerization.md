# Containerise your Golang Microservice

## Considerations

* Target OS
  * local build for Linux / Window / Mac (run binary directly on the developer's machine)
  * build for containerisation, i.e. Linux (run binary inside a Docker container) 
* Building inside a container vs. outside the container

## Docker Multi-Stage Build

Building your Golang Microservice requires a Golang setup on the build machine. You can either provide this setup on the machine itself or move the build completely into a Docker container.

The [`Dockerfile`](../../Dockerfile) provided in this repository uses a so called [multi-stage build](https://docs.docker.com/develop/develop-images/multistage-build/), which means that you have a *build stage* in your Dockerfile that uses a Docker image with Golang tooling to build your binary and later on a stage that uses the artifact built in the build stage to be copied into a Docker image that only provides a bare minimum system.

The benefits of this approach are:

* reproducibility of the build on every machine that runs Docker (i.e. developer workstation and pipeline agent)
* using a lightweight image for the deployable container

## Building the Image

Building the image boils down to

```shell
docker build --no-cache -t golang-microservice .
```

which is wrapped in a convenient `./do` script task, i.e.

```shell
./do.sh build-container
```

## Running the Image

After building the image, you can run it via

```shell
docker run --name golang-microservice -d -p 12345:12345 golang-microservices
```

Building and running is wrapped in another `./do` script task:

```shell
./do.sh run-container
```

## TODOs

- [ ] describe how config is mounted into the container
- [ ] idea of 12-factor app regarding "injecting" the config

## Further resources 
- https://levelup.gitconnected.com/complete-guide-to-create-docker-container-for-your-golang-application-80f3fb59a15e
