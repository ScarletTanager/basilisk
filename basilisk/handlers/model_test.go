package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/basilisk/basilisk/handlers"
)

var _ = Describe("Model", func() {
	var (
		c echo.Context

		// Vars needed for setting up the request
		request        *http.Request
		method, target string
		body           io.Reader
		bodyBytes      []byte

		// For the response
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		if bodyBytes != nil {
			body = bytes.NewReader(bodyBytes)
		} else {
			body = nil
		}
		request = httptest.NewRequest(method, target, body)
		request.Header.Add("Content-type", "application/json")
		c = echo.New().NewContext(request, recorder)
	})

	Describe("CreateModelHandler", func() {
		var (
			rm *handlers.RunningModel
		)

		BeforeEach(func() {
			rm = &handlers.RunningModel{}
			method = http.MethodPost
			target = "/models"
		})

		When("The request body is valid", func() {
			BeforeEach(func() {
				bodyBytes = []byte(`{
					"k": 1
				}`)
			})

			When("The current model is nil", func() {
				BeforeEach(func() {
					rm.Classifier = nil
				})

				It("Returns an echo.HandlerFunc that initializes current", func() {
					h := handlers.CreateModelHandler(rm)
					h(c)
					Expect(rm.Classifier).NotTo(BeNil())
				})
			})

		})

		When("The current model is not nil", func() {

		})
	})
})
