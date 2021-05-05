# golang-microservices

A collection of blog articles on building a microservice in Go

## TODOs

List of topics that we want to cover within the tutorials / blog posts

### Developer Setup

- [ ] Golang setup (how to install & configure)

### Golang Coding

- [ ] architecture: router - controller - services - client - model - storage - helpers / commons
- [ ] simple HTTP router / controller
- [ ] configuration (default config & override from json & read from environment)
- [ ] error handling (mapping errors to HTTP status codes)
- [ ] providing batch commands in separate binaries (DB migration)
- [ ] validation (JSON schema)
- [ ] ORM (sqlboiler...)
- [ ] Feature Toggles
- [ ] authentication & authorization (e.g. use GitHub SSO)
- [ ] different types: requestTypes, domain model, storage types
- [ ] Swagger documentation

### Testing

- [ ] testing containers
- [ ] HTTP integration tests
- [ ] Ginkgo
- [ ] using / generating mocks (tags, `./do` script task...)
- [ ] Contract Tests (consumer tests, provider tests, pact broker, managing pacts in pipeline)
- [ ] Test Coverage

### Build & CI

- [ ] `./do` script (vs pipeline tasks)
- [ ] local dev setup (Docker, compose, config.json)
- [ ] versioning (traceability)
- [ ] linting
- [ ] packaging with Docker (building container images)
- [ ] simple GitHub Actions pipeline

### Deployment

- [ ] helm charts
- [ ] health endpoints
- [ ] configuration & credentials management
