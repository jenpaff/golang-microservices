package api_test

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jenpaff/golang-microservices/api"
	"github.com/jenpaff/golang-microservices/common"
	"github.com/jenpaff/golang-microservices/config"
	custom_errors "github.com/jenpaff/golang-microservices/errors"
	test_helper "github.com/jenpaff/golang-microservices/test-helper"
	"github.com/jenpaff/golang-microservices/users"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("UserController", func() {

	var mockController *gomock.Controller
	var controller *api.Controller
	var router http.Handler
	var userServiceMock *users.MockService

	BeforeEach(func() {
		mockController = gomock.NewController(test_helper.GinkgoTestReporter{})
		userServiceMock = users.NewMockService(mockController)
		controller = api.NewController(config.Config{}, userServiceMock)
		router = api.NewRouter(controller)
	})

	AfterEach(func() {
		mockController.Finish()
	})

	Context("GetUser", func() {

		//TODO: there seems to be a bug with our error handling
		PIt("returns an error when user does not exist", func() {

			userServiceMock.EXPECT().GetUser(gomock.Any(), "invalid-user").Return(nil, fmt.Errorf("error happened: %w", custom_errors.UserNotFound))

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/users/invalid-user", nil)

			router.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
			body := rr.Body.Bytes()
			var error custom_errors.ErrorResponse
			err := json.Unmarshal(body, &error)
			Expect(err).ToNot(HaveOccurred())
			Expect(error.ErrorID).To(Equal("user-1"))
			Expect(error.ErrorMessage).To(Equal("12345"))
		})

		It("returns user when calling /users/{userName}", func() {

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
