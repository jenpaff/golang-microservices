# Testing

Example: we want to write a test for our health Controller

## The Testing Pyramid

- [ ] Describe how we organise our tests
  - [ ] Unit Level
  - [ ] Controller Tests
  - [ ] Integration Tests
  - [ ] End-2-end tests (e.g. testing-container)

## Bootstrapping a Ginkgo Test Suite

Describes how to get started with Ginkgo and generate the required files.

- [ ] describe what each of the file is doing and why we need it

1. From within the `/api` directory run

   ```bash
   ginkgo bootstrap
   ```

   to create a `_suite_test.go` file.

1. For creating a test file, run

   ```bash
   ginkgo generate health-controller.go
   ```

## Ginkgo Test Structure

Ginkgo allowsto structure your tests with `Describe`, `Context`, `It` and `By` to make your tests more expressive.

## Controller Tests

Goal of the controller tests is to test the controller methods for our endpoints. Therefore we run an actual HTTP request against our router and make assertions on the generated response. Since we test an isolated unit of code here (the controller), we count Controller tests as unit tests.

Here are the relevant lines:

```golang
controller := api.NewController()
router := api.NewRouter(controller)
rr := httptest.NewRecorder()
req, _ := http.NewRequest(http.MethodGet, "/health", nil)

router.ServeHTTP(rr, req)

Expect(rr.Code).To(Equal(http.StatusOK))
```

## Mocking

- [ ] describe why we are using mocks
- [ ] describe how to generate mocks (`./do` task)
- [ ] describe how to use mocks in our tests

## Test Task in `./do.sh` Script

For a better developer experience we add a task `./do.sh test` that runs our (unit) tests for us. The task checks whether the Ginkgo binary is installed and if not installs it and afterwards runs the tests.

We use the very same task in our build pipeline to run our tests. Besides reducing the amount of shell snippets in our pipeline definition this also improves "dev-prod-parity" since we can run the very same code that runs the tests in our pipeline on our developer machine and would get fast feedback if something doesn't work as it is expected.

## Generating Test Coverage Reports

Test coverage is a great tool to identify blind spots in our tests. We usually generate a coverage report while running our tests in a pipeline. Some CI servers like Jenkins or Azure DevOps Pipelines allow for showing coverage reports as part of the builds. Unfortunately GitHub actions (on which we based our sample pipeline) so far has no option to show such reports.

We therefore created a task `test-coverage` in our `./do.sh` script to show the test coverage in the browser.

For a detailed explanation of test coverage and associated tooling in Golang check [this blog post](https://blog.golang.org/cover).

- [ ] how to visualise test coverage in a GitHub action (ideally without any 3rd party tools...)

## Further Resources / Thoughts

* [about build tags in integration tests](https://peter.bourgon.org/blog/2021/04/02/dont-use-build-tags-for-integration-tests.html)
* [test coverage in Golang](https://blog.golang.org/cover)


## TODOs

* raise a feature request in Ginkgo that allows this

   ```
   Health Controller service is up
     calling /health returns status up
     /Users/amreivonczettritz/workspace/golang-microservices/api/health_controller_test.go:16
   [It] calling /health returns status up
     /Users/amreivonczettritz/workspace/golang-microservices/api/health_controller_test.go:16
   •
   ------------------------------
   Health Controller service is up
     calling /health returns status 200
     /Users/amreivonczettritz/workspace/golang-microservices/api/health_controller_test.go:32
   [It] calling /health returns status u200
     /Users/amreivonczettritz/workspace/golang-microservices/api/health_controller_test.go:32
   ```

   to be rendered as

   ```
   Health Controller
     service is up
       calling /health returns status up
       calling /health returns status 200
   ```