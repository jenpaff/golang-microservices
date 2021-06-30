//+build integration

package users_test

import (
	"context"
	"database/sql"
	"github.com/jenpaff/golang-microservices/config"
	"github.com/jenpaff/golang-microservices/persistence"
	"github.com/jenpaff/golang-microservices/persistence/models"
	"github.com/jenpaff/golang-microservices/users"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var _ = Describe("Storage", func() {

	ctx := context.Background()
	var sqlDB *sql.DB
	var storage users.Storage
	var persistenceConfig config.PersistenceConfig

	BeforeSuite(func() {
		cfg, err := config.BuildConfig("../config/test.json", "", "")
		Expect(err).ToNot(HaveOccurred())
		persistenceConfig = cfg.Persistence
	})

	BeforeEach(func() {
		sqlDB, _ = persistence.ConnectPostgres(persistenceConfig)
		//Expect(err).ToNot(HaveOccurred())
		storage = users.NewStorage(sqlDB)
	})

	AfterEach(func() {
		_, err := sqlDB.Exec("delete from users")
		Expect(err).ToNot(HaveOccurred())
	})

	It("should create a user and retrieve it successfully", func() {
		name := "User 1"
		user, err := createUserInDb(sqlDB, name, "test@test.com", "1234567")
		Expect(err).ToNot(HaveOccurred())
		storedUser, err := storage.FindByName(ctx, name)
		Expect(err).ToNot(HaveOccurred())
		Expect(storedUser.UserName).To(Equal(user.Username))
		Expect(storedUser.Email).To(Equal(user.Email.String))
		Expect(storedUser.PhoneNumber).To(Equal(user.PhoneNumber.String))
	})

})

func createUserInDb(sqlDB *sql.DB, name, email, phoneNumber string) (models.User, error) {
	ctx := context.Background()
	user := models.User{
		Username:    name,
		Email:       null.StringFrom(email),
		PhoneNumber: null.StringFrom(phoneNumber),
	}
	err := user.Insert(ctx, sqlDB, boil.Infer())
	return user, err
}
