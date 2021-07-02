package validation

import (
	"github.com/golang/mock/gomock"
	test_helper "github.com/jenpaff/golang-microservices/test-helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validation", func() {

	var mockCtrl *gomock.Controller
	type Example struct {
		Name string `json:"name" validate:"notBlank,validRegexInput"`
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(test_helper.GinkgoTestReporter{})
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("field contains invalid characters", func() {
		s := Example{Name: "invalidUserName$$$$"}
		v, err := NewValidate()
		Expect(err).ToNot(HaveOccurred())
		err = v.Struct(s)
		Expect(err).To(HaveOccurred())
	})

	It("field contains only blank spaces", func() {
		s := Example{Name: "    "}
		v, err := NewValidate()
		Expect(err).ToNot(HaveOccurred())
		err = v.Struct(s)
		Expect(err).To(HaveOccurred())
	})

	It("field contains name and blank spaces", func() {
		s := Example{Name: "    test"}
		v, err := NewValidate()
		Expect(err).ToNot(HaveOccurred())
		err = v.Struct(s)
		Expect(err).To(HaveOccurred())
	})

	It("field contains valid username", func() {
		s := Example{Name: "jenpaff123"}
		v, err := NewValidate()
		Expect(err).ToNot(HaveOccurred())
		err = v.Struct(s)
		Expect(err).ToNot(HaveOccurred())
	})
})
