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
	"github.com/ScarletTanager/basilisk/classifiers"
	"github.com/ScarletTanager/basilisk/dsgen"
)

var _ = Describe("DatasetHandlers", func() {
	var (
		err error

		datasetConfig dsgen.DatasetConfig

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
		method = http.MethodPost
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

	Describe("CreateDatasetHandler", func() {
		BeforeEach(func() {
			method = http.MethodPost
			target = "/datasets"

			datasetConfig = dsgen.DatasetConfig{
				RecordCount: 30,
				Classes: map[string][]dsgen.DataSetAttribute{
					"small": {
						{
							Name:                  "length",
							LowerBound:            0.0,
							UpperBound:            10.0,
							AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
						},
						{
							Name:                  "height",
							LowerBound:            0.0,
							UpperBound:            10.0,
							AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
						},
						{
							Name:                  "width",
							LowerBound:            0.0,
							UpperBound:            10.0,
							AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
						},
					},
					"medium": {
						{
							Name:                  "length",
							LowerBound:            10.1,
							UpperBound:            20.0,
							AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
						},
						{
							Name:                  "height",
							LowerBound:            10.1,
							UpperBound:            20.0,
							AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
						},
						{
							Name:                  "width",
							LowerBound:            10.1,
							UpperBound:            20.0,
							AllocationsByQuintile: []float64{20.0, 20.0, 20.0, 20.0, 20.0},
						},
					},
				},
			}
		})

		When("The configuration is valid", func() {
			BeforeEach(func() {
				bodyBytes, err = json.Marshal(datasetConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(bodyBytes).NotTo(BeNil())
			})

			It("Returns a 200", func() {
				handlers.CreateDatasetHandler(c)
				Expect(recorder.Result().StatusCode).To(Equal(http.StatusOK))
			})

			It("Returns a dataset of the correct size", func() {
				handlers.CreateDatasetHandler(c)
				respBytes, err := io.ReadAll(recorder.Result().Body)
				Expect(err).NotTo(HaveOccurred())

				var dataset classifiers.DataSet
				err = json.Unmarshal(respBytes, &dataset)
				Expect(err).NotTo(HaveOccurred())

				Expect(dataset.Records).To(HaveLen(datasetConfig.RecordCount))
			})

			It("Returns a dataset containing the correct classes", func() {
				handlers.CreateDatasetHandler(c)
				respBytes, err := io.ReadAll(recorder.Result().Body)
				Expect(err).NotTo(HaveOccurred())

				var dataset classifiers.DataSet
				err = json.Unmarshal(respBytes, &dataset)
				Expect(err).NotTo(HaveOccurred())

				Expect(dataset.ClassNames).To(ConsistOf("small", "medium"))
			})

			It("Returns a dataset containing the correct attributes", func() {
				handlers.CreateDatasetHandler(c)
				respBytes, err := io.ReadAll(recorder.Result().Body)
				Expect(err).NotTo(HaveOccurred())

				var dataset classifiers.DataSet
				err = json.Unmarshal(respBytes, &dataset)
				Expect(err).NotTo(HaveOccurred())

				Expect(dataset.AttributeNames).To(ConsistOf("length", "height", "width"))
			})
		})

		When("The configuration is not a valid JSONification of a DatasetConfig", func() {
			BeforeEach(func() {
				bodyBytes = []byte(`doesthis { look} like valid JSONto you?`)
			})

			It("Returns a 400", func() {
				handlers.CreateDatasetHandler(c)
				Expect(recorder.Result().StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		When("The configuration is not valid for Dataset creation", func() {
			BeforeEach(func() {
				// Remove one attribute from one class
				datasetConfig.Classes["medium"] = datasetConfig.Classes["medium"][:2]
			})

			It("Returns a 400", func() {
				handlers.CreateDatasetHandler(c)
				Expect(recorder.Result().StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
