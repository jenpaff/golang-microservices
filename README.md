# golang-microservices

A collection of blog articles on building a microservice in Go

## TODOs

List of topics that we want to cover within the tutorials / blog posts

## TO FIX
- [ ] Database connection

### Developer Setup

- [ ] Golang setup (how to install & configure)

### Golang Coding

- [ ] architecture: router - controller - services - client - model - storage - helpers / commons => at the end! 
- [ ] simple HTTP router / controller / go-restful 
- [X] configuration (default config & override from json & read from environment) => Jen
- [X] error handling (mapping errors to HTTP status codes) => Michael
- [X] providing batch commands in separate binaries (DB migration)
- [X] validation (JSON schema) => Jen
- [x] ORM (sqlboiler...)
- [X] Feature Toggles => Jen
- [ ] authentication & authorization (e.g. use GitHub SSO) => we don't need the codebase for this
- [ ] different types: requestTypes, domain model, storage types
- [X] Swagger documentation => Jen
- [ ] Performance & Benchmarking

### Testing

- [ ] e2e testing with testing-containers (health + endpoint that checks db, ...)
- [ ] HTTP integration tests (all endpoints)
- [ ] Ginkgo (parameterized tests)
- [ ] using / generating mocks (gomock, go generate tags, `./do` script task...)
- [ ] Contract Tests (consumer tests, provider tests, pact broker, managing pacts in pipeline) => Amrei
- [ ] Test Coverage (check out ginkgo test coverage - how to check integration & contract tests)  => Michael

### Build & CI

- [X] `./do` script (vs pipeline tasks)
- [X] local dev setup (Docker, compose, config.json)
- [ ] versioning (traceability)
- [ ] linting
- [X] packaging with Docker (`Dockerfile`, building container images)
- [ ] simple GitHub Actions pipeline
- [ ] using go modules from private repositories (`go.customer.com/...`) with proper `.gitconfig`...
- [x] provide a `Dockerfile` example for a multi-stage build

### Deployment

- [ ] helm charts
- [ ] health endpoints
- [ ] configuration & credentials management

### Documentation

- [ ] How to write a proper `README`

### Anti-Patterns

- [ ] shared libraries

## References

* [GitHub Actions for Golang](https://github.com/mvdan/github-actions-golang)
* [Set up Golang on Mac with Homebrew](https://jimkang.medium.com/install-go-on-mac-with-homebrew-5fa421fc55f5)
* [go-restful blog post](http://ernestmicklei.com/2012/11/go-restful-first-working-example/)
