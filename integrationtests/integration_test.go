//+build integration

package integrationtests

import (
	"encoding/json"
	"github.com/jenpaff/golang-microservices/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
)

var _ = Describe("Golang Service", func() {

	var golangService GolangService

	BeforeEach(func() {
		golangService = NewGolangService()
		golangService.Start()
	})

	AfterEach(func() {
		golangService.Stop()
	})

	Context("when it is started", func() {

		It("should have health endpoint return status ok", func() {
			By("By returning a 200 status code")
			response := golangService.Get("/health", map[string]string{})
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			By("By having a valid json body")
			bodyBytes, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())
			healthResponse := api.Health{}
			err = json.Unmarshal(bodyBytes, &healthResponse)
			Expect(err).ToNot(HaveOccurred())
			By("By having the correct name and status up")
			Expect(healthResponse.Status).To(Equal("up"))
		})
	})
})
