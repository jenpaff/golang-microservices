package featuretoggles

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/golang/mock/gomock"
	"github.com/jenpaff/golang-microservices/config"
	"github.com/jenpaff/golang-microservices/test-helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
)

var _ = Describe("Feature Toggles", func() {

	var mockCtrl *gomock.Controller

	BeforeEach(func() {
		mockCtrl = gomock.NewController(test_helper.GinkgoTestReporter{})
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("IsEnabled", func() {
		It("should return the toggle value from Config if Request does not contain it", func() {
			appConfig := config.Config{
				FeatureToggles: map[string]bool{
					"toggle1": true,
					"toggle2": false,
				},
			}

			ft := NewFeatureToggles(&appConfig, createFakeRequest(""))

			Expect(ft.IsEnabled("toggle1")).To(BeTrue())
			Expect(ft.IsEnabled("toggle2")).To(BeFalse())
		})

		It("should return false if the toggle is not defined", func() {
			appConfig := config.Config{
				FeatureToggles: map[string]bool{
					"toggle1": true,
				},
			}

			ft := NewFeatureToggles(&appConfig, createFakeRequest(""))

			Expect(ft.IsEnabled("undefinedToggle")).To(BeFalse())
		})

		It("should override the toggle value based on Request query params", func() {
			appConfig := config.Config{
				FeatureToggles: map[string]bool{
					"toggle1": false,
				},
			}

			ft := NewFeatureToggles(&appConfig, createFakeRequest("toggle1=true"))

			Expect(ft.IsEnabled("toggle1")).To(BeTrue())
		})

	})
})

func createFakeRequest(queryParams string) *restful.Request {
	httpRequest := httptest.NewRequest("GET", "http://www.test.com?"+queryParams, nil)
	request := restful.NewRequest(httpRequest)

	return request
}
