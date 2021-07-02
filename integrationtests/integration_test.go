//+build integration

package integrationtests

import (
	"bytes"
	"encoding/json"
	"github.com/jenpaff/golang-microservices/api"
	"github.com/jenpaff/golang-microservices/common"
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

	Context("Users API", func() {
		Context("with new feature", func() {
			It("should create and retrieve a given user", func() {
				By("By returning a 200 status code when creating a user")

				user, err := json.Marshal(&api.UserCreationRequest{
					UserName:    "jenpaff1",
					Email:       "jenpaff1@test.com",
					PhoneNumber: "0123456781",
				})
				Expect(err).ToNot(HaveOccurred())
				response := golangService.Post("/users?enableNewFeature=true", map[string]string{}, bytes.NewReader(user))
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				By("By returning a 200 status code when retrieving a user")
				response = golangService.Get("/users/jenpaff1", map[string]string{})
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				By("By having a valid json body")
				bodyBytes, err := ioutil.ReadAll(response.Body)
				Expect(err).ToNot(HaveOccurred())
				userResponse := common.User{}
				err = json.Unmarshal(bodyBytes, &userResponse)
				Expect(err).ToNot(HaveOccurred())
				By("By having the correct name and status up")
				Expect(userResponse.UserName).To(Equal("jenpaff1"))
			})
		})

		Context("without new feature", func() {
			It("should create and retrieve a given user", func() {
				By("By returning a 200 status code when creating a user")

				user, err := json.Marshal(&api.UserCreationRequest{
					UserName:    "jenpaff",
					Email:       "jenpaff@test.com",
					PhoneNumber: "012345678",
				})
				Expect(err).ToNot(HaveOccurred())
				response := golangService.Post("/users", map[string]string{}, bytes.NewReader(user))
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				By("By returning a 200 status code when retrieving a user")
				response = golangService.Get("/users/jenpaff", map[string]string{})
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				By("By having a valid json body")
				bodyBytes, err := ioutil.ReadAll(response.Body)
				Expect(err).ToNot(HaveOccurred())
				userResponse := common.User{}
				err = json.Unmarshal(bodyBytes, &userResponse)
				Expect(err).ToNot(HaveOccurred())
				By("By having the correct name and status up")
				Expect(userResponse.UserName).To(Equal("jenpaff"))
			})
		})
	})
})
