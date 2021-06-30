//+build unit

package api_test

import (
	"encoding/json"
	"github.com/jenpaff/golang-microservices/api"
	"github.com/jenpaff/golang-microservices/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Health Controller", func() {

	Context("service is up", func() {

		It("calling /health returns status up", func() {
			controller := api.NewController(config.Config{
				Name: "Golang Test Service",
			})
			router := api.NewRouter(controller)
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/health", nil)

			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			body := rr.Body.Bytes()
			var health api.Health
			err := json.Unmarshal(body, &health)
			Expect(err).ToNot(HaveOccurred())
			Expect(health.Status).To(Equal("up"))
		})
	})
})
