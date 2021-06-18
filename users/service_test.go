//+build unit

package users_test

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jenpaff/golang-microservices/errors"
	"github.com/jenpaff/golang-microservices/persistence/models"
	test_helper "github.com/jenpaff/golang-microservices/test-helper"
	"github.com/jenpaff/golang-microservices/users"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/volatiletech/null/v8"
)

var _ = Describe("Service", func() {

	var userService users.Service
	var mockCtrl *gomock.Controller
	var mockStorage *users.MockStorage

	ctx := context.Background()

	BeforeEach(func() {
		mockCtrl = gomock.NewController(test_helper.GinkgoTestReporter{})
		mockStorage = users.NewMockStorage(mockCtrl)
		userService = users.NewService(mockStorage)
	})

	Context("Get user", func() {

		username := "test"

		It("will return error if user can not be found", func() {
			mockStorage.EXPECT().FindByName(gomock.Any(), "test").Return(nil, fmt.Errorf("user was not found %w", errors.UserNotFound))

			_, err := userService.GetUser(ctx, username)

			Expect(err).To(HaveOccurred())
		})

		It("will successfully retrieve a user by username", func() {

			email := "test@test.com"
			phoneNumber := "test@test.com"
			mockStorage.EXPECT().FindByName(gomock.Any(), "test").Return(&models.User{
				ID:          123,
				Username:    username,
				Email:       null.StringFrom(email),
				PhoneNumber: null.StringFrom(phoneNumber),
			}, nil)

			returnedUser, err := userService.GetUser(ctx, username)

			Expect(err).ToNot(HaveOccurred())
			Expect(returnedUser.Username).To(Equal(username))
		})

	})

})
