//+build unit

package errors_test

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/jenpaff/golang-microservices/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("ErrorHandler", func() {

	DescribeTable("handles error", func(err error, expectedStatus int) {
		rr := httptest.NewRecorder()

		errorHandler := errors.ErrorHandler(func(_ *restful.Request, _ *restful.Response) error {
			return err
		})

		res := restful.NewResponse(rr)
		res.SetRequestAccepts(restful.MIME_JSON)

		errorHandler(nil, res)

		Expect(rr.Code).To(Equal(expectedStatus))
	},
		Entry("returns expected http status given a known error", fmt.Errorf("error happened: %w", errors.UserNotFound), http.StatusNotFound),
		Entry("returns internal server error given a unknown error", fmt.Errorf("error happened"), http.StatusInternalServerError),
	)
})
