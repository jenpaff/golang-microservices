//+build container

package containertests

import (
	"bytes"
	"context"
	"encoding/json"
	docker_types "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/jenpaff/golang-microservices/api"
	"github.com/jenpaff/golang-microservices/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
)

var _ = Describe("Golang Service", func() {

	imageName := "golang-microservices"
	var g GolangService

	BeforeSuite(func() {
		var err error
		g, err = NewGolangService(imageName, context.Background())
		Expect(err).ToNot(HaveOccurred())
		g.Start()
	})

	AfterSuite(func() {
		g.Stop()
		removeContainer(imageName)
		removeContainer("golang-microservices")
	})

	AfterEach(func() {
		err := CleanUpDatabase()
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("The Golang Service", func() {
		Context("when it is started", func() {
			It("should have health endpoint return status ok", func() {
				By("By returning a 200 status code")
				response := g.Get("/health", map[string]string{})
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				By("By having a valid json body")
				bodyBytes, err := ioutil.ReadAll(response.Body)
				Expect(err).ToNot(HaveOccurred())

				healthResponse := api.Health{}
				err = json.Unmarshal(bodyBytes, &healthResponse)
				Expect(err).ToNot(HaveOccurred())

				By("By having the correct name and status up")
				Expect(healthResponse.Name).To(Equal("Golang Service"))
				Expect(healthResponse.Status).To(Equal("up"))
			})
		})

		Context("Users API", func() {
			Context("without new feature", func() {
				It("should create and retrieve a given user", func() {
					By("By returning a 200 status code when creating a user")

					user, err := json.Marshal(&api.UserCreationRequest{
						UserName:    "jenpaff",
						Email:       "jenpaff@test.com",
						PhoneNumber: "012345678",
					})
					Expect(err).ToNot(HaveOccurred())
					response := g.Post("/users?enableNewFeature=false", map[string]string{}, bytes.NewReader(user))
					Expect(response.StatusCode).To(Equal(http.StatusOK))
					By("By returning a 200 status code when retrieving a user")
					response = g.Get("/users/jenpaff", map[string]string{})
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

			Context("with new feature", func() {
				It("should create and retrieve a given user", func() {
					By("By returning a 200 status code when creating a user")

					user, err := json.Marshal(&api.UserCreationRequest{
						UserName:    "jenpaff",
						Email:       "jenpaff@test.com",
						PhoneNumber: "012345678",
					})
					Expect(err).ToNot(HaveOccurred())
					response := g.Post("/users", map[string]string{}, bytes.NewReader(user))
					Expect(response.StatusCode).To(Equal(http.StatusOK))
					By("By returning a 200 status code when retrieving a user")
					response = g.Get("/users/jenpaff", map[string]string{})
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
})

func removeContainer(imageName string) {
	c, err := client.NewClientWithOpts()
	Expect(err).ToNot(HaveOccurred())

	filterArgs := filters.NewArgs()
	filterArgs.Add("ancestor", imageName)
	containers, err := c.ContainerList(context.Background(), docker_types.ContainerListOptions{All: true, Filters: filterArgs})
	Expect(err).ToNot(HaveOccurred())

	for _, container := range containers {
		err = c.ContainerRemove(context.Background(), container.ID, docker_types.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
		Expect(err).ToNot(HaveOccurred())
	}
}
