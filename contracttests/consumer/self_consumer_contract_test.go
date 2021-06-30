//+build unit

package consumer_test

import (
	"fmt"
	"github.com/jenpaff/golang-microservices/contracttests/consumer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"

	"github.com/pact-foundation/pact-go/dsl"
	"net/http"
)

var _ = Describe("SelfConsumerContract", func() {

	var pact *dsl.Pact

	BeforeSuite(func() {
		pact = setupPact()
		pact.Setup(true)
	})

	AfterSuite(func() {
		_ = pact.WritePact()
		pact.Teardown()
	})

	It("returns a 404 if user does not exist", func() {
		pact.
			AddInteraction().
			Given("user1 does not exist").
			UponReceiving("A request to fetch user1").
			WithRequest(dsl.Request{
				Method: http.MethodGet,
				Path:   dsl.String("/users/user1"),
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.Regex("application/json; charset=utf-8", "application/json.*"),
				},
			}).
			WillRespondWith(dsl.Response{
				Status: http.StatusNotFound,
				// TODO: errorresponse
			})
		userClient := consumer.UserClient{
			BaseURL: fmt.Sprintf("http://localhost:%d", pact.Server.Port),
		}
		var test = func() error {
			returnedUser, err := userClient.GetUser("user1")
			Expect(returnedUser).To(BeNil())
			Expect(err).To(HaveOccurred())
			return nil
		}

		if err := pact.Verify(test); err != nil {
			Expect(err).ToNot(HaveOccurred(), "Error on Verify")
		}
	})

})

func setupPact() *dsl.Pact {
	pactsDirAbs, _ := filepath.Abs("../pacts/")
	pact := &dsl.Pact{
		Consumer:          "My external Service",
		Provider:          "Golang Service",
		Host:              "localhost",
		PactFileWriteMode: "merge",
		PactDir:           pactsDirAbs,
		LogLevel:          "WARN",
		DisableToolValidityCheck: false,
	}
	return pact
}
