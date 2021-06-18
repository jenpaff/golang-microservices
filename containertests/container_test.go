//+build container

package containertests

import (
	"context"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
)

var _ = Describe("Golang Service", func() {

	imageName := "golang-service:local-dev-version"
	var p GolangService

	BeforeSuite(func() {
		var err error
		p, err = NewGolangService(imageName, context.Background())
		Expect(err).ToNot(HaveOccurred())
		p.Start()
	})

	AfterSuite(func() {
		p.Stop()
		removeContainer(imageName)
		removeContainer("golang-service")
	})

	AfterEach(func() {
		err := CleanUpDatabase()
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("The Golang Service", func() {
		Context("when it is started", func() {
			It("should have health endpoint return status ok", func() {
				By("By returning a 200 status code")
				response := p.Get("/health", map[string]string{})
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				By("By having a valid json body")
				bodyBytes, err := ioutil.ReadAll(response.Body)
				Expect(err).ToNot(HaveOccurred())

				healthResponse := health.ServiceInfo{}
				err = json.Unmarshal(bodyBytes, &healthResponse)
				Expect(err).ToNot(HaveOccurred())

				By("By having the correct name and status up")
				Expect(healthResponse.Name).To(Equal("Golang Service"))
				Expect(healthResponse.Status).To(Equal("up"))
				Expect(healthResponse.DatabaseAvailable).To(BeTrue())
			})
		})
	})
})
