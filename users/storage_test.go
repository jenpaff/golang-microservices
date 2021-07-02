//+build integration

package users_test

import (
	"context"
	"database/sql"
	"github.com/jenpaff/golang-microservices/config"
	"github.com/jenpaff/golang-microservices/persistence"
	"github.com/jenpaff/golang-microservices/users"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage", func() {

	ctx := context.Background()
	var sqlDB *sql.DB
	var storage users.Storage
	var err error
	var cfg config.Config

	BeforeSuite(func() {
		cfg, err = config.BuildConfig("../config/test.json", "", "")
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		sqlDB, err = persistence.ConnectPostgres(cfg.Persistence)
		Expect(err).ToNot(HaveOccurred())
		storage = users.NewStorage(sqlDB)
	})

	AfterEach(func() {
		_, err = sqlDB.Exec("delete from users")
		Expect(err).ToNot(HaveOccurred())
	})

	It("should create a user and retrieve it successfully", func() {
		name := "User 1"
		email := "test@test.com"
		phone := "1234567"
		_, err := storage.Create(ctx, name, email, phone)
		Expect(err).ToNot(HaveOccurred())
		storedUser, err := storage.FindByName(ctx, name)
		Expect(err).ToNot(HaveOccurred())
		Expect(storedUser.UserName).To(Equal(name))
		Expect(storedUser.Email).To(Equal(email))
		Expect(storedUser.PhoneNumber).To(Equal(phone))
	})

})
