//go:generate mockgen -destination=feature_toggles_mock.go -package=featuretoggles github.com/jenpaff/golang-microservices/featuretoggles FeatureToggles

package featuretoggles

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/go-playground/log"
	"github.com/jenpaff/golang-microservices/config"
	"strconv"
)

type FeatureToggles interface {
	IsEnabled(toggleName string) bool
}

type featureToggles struct {
	appConfig   *config.Config
	httpRequest *restful.Request
}

func NewFeatureToggles(appConfig *config.Config, httpRequest *restful.Request) FeatureToggles {
	return &featureToggles{appConfig: appConfig, httpRequest: httpRequest}
}

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
