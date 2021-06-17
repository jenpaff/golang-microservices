# Feature Toggles

## Motivation
Generally our team advocates trunk-based development and continuous integration & deployment. While aiming at deploying to PROD with every commit, feature toggles enable us to control which features will be exposed. 

## Code

Let's say we introduce a new feature which we want to enable on our dev environment first without enabling it on our PROD environment. 
Our dev configuration looks like this: 

```json
{
  "featureToggles": {
    "enableNewFeature": true
  }
}

```

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
