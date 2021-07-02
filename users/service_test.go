//+build unit

package users_test

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jenpaff/golang-microservices/common"
	"github.com/jenpaff/golang-microservices/errors"
	test_helper "github.com/jenpaff/golang-microservices/test-helper"
	"github.com/jenpaff/golang-microservices/users"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			phoneNumber := "12345"
			mockStorage.EXPECT().FindByName(gomock.Any(), "test").Return(&common.User{
				UserName:    username,
				Email:       email,
				PhoneNumber: phoneNumber,
			}, nil)

			returnedUser, err := userService.GetUser(ctx, username)

			Expect(err).ToNot(HaveOccurred())
			Expect(returnedUser.UserName).To(Equal(username))
		})

	})

	Context("Create user", func() {

		It("will return error if user can not be found", func() {

			username := "test"
			email := "test@test.com"
			phone := "1234"

			mockStorage.EXPECT().Create(gomock.Any(), username, "test@test.com", "1234").Return(nil, fmt.Errorf("error storing user %w", errors.DatabaseError))

			_, err := userService.CreateUser(ctx, username, email, phone)

			Expect(err).To(HaveOccurred())
		})

		It("will successfully create a user", func() {

			username := "test"
			email := "test@test.com"
			phone := "1234"

			mockStorage.EXPECT().Create(gomock.Any(), username, email, phone).Return(&common.User{
				UserName:    username,
				Email:       email,
				PhoneNumber: phone,
			}, nil)

			returnedUser, err := userService.CreateUser(ctx, username, email, phone)

			Expect(err).ToNot(HaveOccurred())
			Expect(returnedUser.UserName).To(Equal(username))
		})

	})
})
