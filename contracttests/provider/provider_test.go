// +build unit

package provider_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jenpaff/golang-microservices/api"
	"github.com/jenpaff/golang-microservices/config"
	"github.com/jenpaff/golang-microservices/users"
	"github.com/jenpaff/golang-microservices/validation"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

const (
	runAllLocalTests = "RUN_ALL_LOCAL_TESTS"
	pactLocalFiles   = "PACT_LOCAL_FILES"
	pactUserName     = "PACT_USER"
	pactPassword     = "PACT_PASSWORD"
	pactPublishFlag  = "PACT_PUBLISH_FLAG"
	providerVersion  = "PROVIDER_VERSION"
	localContracts   = "../localcontracts"
)

func PTestProvider(t *testing.T) {
	handler, finalizer := setup(t)
	defer finalizer()

	pact := &dsl.Pact{
		Provider: "Golang Service",
		LogLevel: "WARN",
	}
	defer pact.Teardown()

	go startServer(t, handler)

	publishPacts := false
	providerVersionValue := "0.0.0"
	if os.Getenv(pactPublishFlag) == "true" {
		publishPacts = true
		providerVersionValue = os.Getenv(providerVersion)
		assert.NotEmpty(t, providerVersionValue)
	}

	brokerUrl := "https://golang-pact-broker.aws.net"
	brokerUser := os.Getenv(pactUserName)
	brokerPass := os.Getenv(pactPassword)
	var localPacts []string
	if localPactsString := os.Getenv(pactLocalFiles); localPactsString != "" {
		if !(filepath.IsAbs(localPactsString)) {
			localPactsString = filepath.Join("..", localPactsString)
		}
		localPacts = []string{localPactsString}
		// and disable broker pacts
		brokerUrl = ""
		brokerUser = ""
		brokerPass = ""
	}

	if os.Getenv(runAllLocalTests) == "true" {
		localTestFiles, err := getLocalTestFiles()
		assert.NoError(t, err)
		localPacts = localTestFiles

		brokerUrl = ""
		brokerUser = ""
		brokerPass = ""
	}

	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL:            "http://localhost:12345",
		BrokerURL:                  brokerUrl,
		BrokerUsername:             brokerUser,
		BrokerPassword:             brokerPass,
		PactURLs:                   localPacts,
		PublishVerificationResults: publishPacts,
		FailIfNoPactsFound:         true,
		ProviderVersion:            providerVersionValue,
		CustomProviderHeaders:      []string{fmt.Sprintf("Authorization: Bearer %s", "Bearer 1234")},
		StateHandlers: types.StateHandlers{
			"Statehandler text": func() error {
				return nil
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
	}
}

func startServer(t *testing.T, handler http.Handler) {
	assert.NoError(t, http.ListenAndServe("localhost:12345", handler))
}

func setup(t *testing.T) (handler http.Handler, finalizer func()) {
	fmt.Println("Setting up")
	mockController := gomock.NewController(t)
	userServiceMock := users.NewMockService(mockController)
	validator, _ := validation.NewValidate()
	controller := api.NewController(config.Config{}, userServiceMock, validator)
	router := api.NewRouter(controller)

	return router, func() {
		mockController.Finish()
	}
}

func getLocalTestFiles() ([]string, error) {
	var filenames []string
	contracts, err := ioutil.ReadDir(localContracts)
	if err != nil {
		return nil, err
	}

	for _, contract := range contracts {
		localPactsString := contract.Name()
		filenames = append(filenames, fmt.Sprintf("%s/%s", localContracts, localPactsString))
	}

	return filenames, nil
}
