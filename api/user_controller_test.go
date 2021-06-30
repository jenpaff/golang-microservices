package api_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/jenpaff/golang-microservices/api"
	"github.com/jenpaff/golang-microservices/common"
	"github.com/jenpaff/golang-microservices/config"
	test_helper "github.com/jenpaff/golang-microservices/test-helper"
	"github.com/jenpaff/golang-microservices/users"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("UserController", func() {

	var mockController *gomock.Controller
	var router http.Handler
	var userServiceMock *users.MockService
	var controller *api.Controller

	BeforeEach(func() {
		mockController = gomock.NewController(test_helper.GinkgoTestReporter{})
		controller = api.NewController(config.Config{})
		userServiceMock = users.NewMockService(mockController)
		router = api.NewRouter(controller)
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Context("GetUser", func() {

		PIt("calling /users/{userName} returns user", func() {

			userServiceMock.EXPECT().GetUser(gomock.Any(), "user-id-1").Return(&common.User{
				UserName:    "user-1",
				PhoneNumber: "12345",
				Email:       "test@test.com",
			}, nil)

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/users/user-id-1", nil)

			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			body := rr.Body.Bytes()
			var user common.User
			err := json.Unmarshal(body, &user)
			Expect(err).ToNot(HaveOccurred())
			Expect(user.UserName).To(Equal("user-1"))
			Expect(user.PhoneNumber).To(Equal("12345"))
			Expect(user.Email).To(Equal("test@test.com"))
		})
	})
})
