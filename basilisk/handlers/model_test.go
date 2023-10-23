package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/basilisk/basilisk/handlers"
	"github.com/ScarletTanager/basilisk/classifiers/knn"
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
			rm *handlers.RunningModels
		)

		BeforeEach(func() {
			rm = &handlers.RunningModels{}
			method = http.MethodPost
			target = "/models"
		})

		When("The request body is valid", func() {
			BeforeEach(func() {
				bodyBytes = []byte(`{
					"k": 1
				}`)
			})

			When("No models exist", func() {
				BeforeEach(func() {
					rm.Classifiers = nil
				})

				It("Returns an echo.HandlerFunc that creates a new model", func() {
					h := handlers.CreateModelHandler(rm)
					h(c)
					Expect(rm.Classifiers).NotTo(BeNil())
				})

				It("Returns an HTTP 200", func() {
					h := handlers.CreateModelHandler(rm)
					h(c)
					resp := recorder.Result()
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
				})

				It("Returns the model metadata", func() {
					h := handlers.CreateModelHandler(rm)
					h(c)
					resp := recorder.Result()
					body, _ := io.ReadAll(resp.Body)
					m := &handlers.Model{}
					Expect(json.Unmarshal(body, m)).NotTo(HaveOccurred())
					Expect(m.ID).To(Equal(0))
					Expect(m.K).To(Equal(1))
				})
			})

			When("Models exist", func() {
				var (
					knnc *knn.KNearestNeighborClassifier
				)

				BeforeEach(func() {
					knnc, _ = knn.New(1)
					_, err := rm.Add(knnc)
					Expect(err).NotTo(HaveOccurred())
					Expect(rm.Classifiers).To(HaveLen(1))
				})

				It("Adds another model", func() {
					h := handlers.CreateModelHandler(rm)
					h(c)
					Expect(rm.Classifiers).To(HaveLen(2))
				})

				It("Returns a status OK", func() {
					h := handlers.CreateModelHandler(rm)
					h(c)
					resp := recorder.Result()
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
				})

				It("Returns the model metadata", func() {
					h := handlers.CreateModelHandler(rm)
					h(c)
					resp := recorder.Result()
					body, _ := io.ReadAll(resp.Body)
					m := &handlers.Model{}
					Expect(json.Unmarshal(body, m)).NotTo(HaveOccurred())
					Expect(m.ID).To(Equal(1))
					Expect(m.K).To(Equal(1))
				})
			})
		})

		When("The request body is invalid", func() {
			BeforeEach(func() {
				bodyBytes = []byte(`{
					hooboah
				}`)
			})

			It("Returns an HTTP 400", func() {
				h := handlers.CreateModelHandler(rm)
				h(c)
				resp := recorder.Result()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
