# Feature Toggles

## Motivation
Generally our team advocates trunk-based development and continuous integration & deployment. 
While aiming at deploying to PROD with every commit, feature toggles enable us to control which features will be exposed. 
One important thing to note is that generally we want our feature toggles to be short living.

## Usage 

### Using the toggle

Let's say we introduce a new feature which we want to enable on our dev environment first without enabling it on our PROD environment.
Our dev configuration looks like this:

```json
{
  "featureToggles": {
    "enableNewFeature": true
  }
}
```

In our code we can now use our featureToggleService to fetch the toggle via our config and use the toggle as we please: 

```
ft := featuretoggles.NewFeatureToggles(c.config, request)
if ft.IsEnabled("enableNewFeature") {
    // do this
} else {
    // do that 
}
```

For ease of use we've also enabled a toggle override via query parameter e.g. 
```bash
curl -X POST http://www.test.com/users?enableNewFeature=true
- d '{
  "username": "test",
  "email": "test@test.com"
}'
```

### Implementing the toggle

We've had multiple conversations about how to deal with duplicated code, tests for feature toggle combinations etc. 
Let me walk through our process by example 

**Step 1:** Implement new toggle in our dev config file
```json
{
  "featureToggles": {
    "enableNewFeature": true
  }
}
```
**Step 2:** duplicate "the old implementation and related tests" such that we do not make any changes to the current implementation 
```go
package users

type Service interface {
	CreateUser(ctx context.Context, userName, email, phoneNumber string) (*common.User, error)
	CreateUserWithNewFeature(ctx context.Context, userName, email, phoneNumber string) (*common.User, error)
}
```

**Step 3:** use the feature toggle to switch between old and new implementation
```go
var createdUser *common.User

ft := featuretoggles.NewFeatureToggles(&c.Cfg, req)
if ft.IsEnabled("enableNewFeature") {
    createdUser, err = c.userService.CreateUserWithNewFeature(req.Request.Context(), creationRequest.UserName, creationRequest.Email, creationRequest.PhoneNumber)
    if err != nil {
        return fmt.Errorf("could not create user: %w", err)
    }
} else {
    createdUser, err = c.userService.CreateUser(req.Request.Context(), creationRequest.UserName, creationRequest.Email, creationRequest.PhoneNumber)
    if err != nil {
        return fmt.Errorf("could not create user: %w", err)
    }
}
```

Since we're using ginkgo in our tests, we can easily 

```go
Context("CreateUser", func() {

    Context("without new feature", func() {
        // tests without feature
    })

    Context("with new feature", func() {
        // tests with feature
})
```

I know it feels strange with much duplication code, however, it helps us to easily clean up afterwards.

**Step 4:** Having tested our new code on all environments, we usually have a separate clean up task to: 
1. remove old implementation code and tests
2. remove the toggle and update tests
3. make sure all tests run correctly
4. rename / refactoring

## Code

### How the feature toggle service is implemented

Let's look at our `featureToggles` struct which holds two values namely a configuration and a request (in our specific case we are using `go-restful` and therefore it's a `*restful.Request`).  

```go
type featureToggles struct {
    appConfig   *config.Config
    httpRequest *restful.Request
}
```

This is useful so we can accept a feature toggle both via config file as well as via request parameter as an override.
To verfiy whether a certain feature toggle is enabled we 

```go
func (ft *featureToggles) IsEnabled(toggleName string) bool {
	toggleState := ft.appConfig.FeatureToggles[toggleName]

	toggleOverride := ft.httpRequest.QueryParameters(toggleName)
	if len(toggleOverride) > 0 {
		log.Infof("overriding toggle '%v' from request - switching from '%v' to '%v'", toggleName, toggleState, toggleOverride[0])
		toggleState, _ = strconv.ParseBool(toggleOverride[0])
	}

	log.Infof("toggle '%v' state set to '%v'", toggleName, toggleState)
	return toggleState
}
```


## Further Resources
* [Managing feature toggles in teams](https://www.thoughtworks.com/insights/blog/managing-feature-toggles-teams)
