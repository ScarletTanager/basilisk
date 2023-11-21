package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ScarletTanager/basilisk/basilisk/handlers"
	"github.com/ScarletTanager/basilisk/basilisk/model"
	"github.com/ScarletTanager/basilisk/classifiers/knn"
)

var _ = Describe("Model", func() {
	var (
		c    echo.Context
		rm   *model.RunningModels
		knnc *knn.KNearestNeighborClassifier

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
		rm = &model.RunningModels{}
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
		BeforeEach(func() {
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
					m := &model.Model{}
					Expect(json.Unmarshal(body, m)).NotTo(HaveOccurred())
					Expect(m.ID).To(Equal(0))
					Expect(m.K).To(Equal(1))
				})
			})

			When("Models exist", func() {
				BeforeEach(func() {
					knnc, _ = knn.New(1, "")
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
					m := &model.Model{}
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

	Describe("TrainModelHandler", func() {
		BeforeEach(func() {
			target = "/models/0/trainingdata"
			method = http.MethodPut
			knnc, _ = knn.New(1, "")
		})

		When("The model has been set in the context", func() {
			JustBeforeEach(func() {
				c.Set(handlers.ContextKeyModel, knnc)
			})

			When("The body is valid", func() {
				BeforeEach(func() {
					bodyBytes, _ = os.ReadFile("../../datasets/shorebirds.json")
				})

				When("The model is untrained", func() {
					BeforeEach(func() {
						Expect(knnc.TrainingData).To(BeNil())
					})

					It("Returns an echo.HandlerFunc which trains the correct classifier", func() {
						h := handlers.TrainModelHandler(rm)
						h(c)
						Expect(knnc.TrainingData).NotTo(BeNil())
					})

					It("Returns an echo.HandlerFunc which returns a 200", func() {
						h := handlers.TrainModelHandler(rm)
						h(c)
						resp := recorder.Result()
						Expect(resp.StatusCode).To(Equal(http.StatusOK))
					})
				})

				When("The model has previously been trained", func() {
					Context("With the same data", func() {
						JustBeforeEach(func() {
							handlers.TrainModelHandler(rm)(c)
						})

						It("Retrains the model", func() {
							orig := knnc.TrainingData.Records

							request = httptest.NewRequest(method, target, body)
							request.Header.Add("Content-type", "application/json")
							newCtx := echo.New().NewContext(request, &httptest.ResponseRecorder{})
							newCtx.Set(handlers.ContextKeyModel, knnc)

							handlers.TrainModelHandler(rm)(newCtx)
							Expect(knnc.TrainingData.Records).NotTo(ConsistOf(orig))
						})
					})

					// Context("With different data", func() {
					// 	JustBeforeEach(func() {
					// 		handlers.TrainModelHandler(rm)(c)
					// 	})

					// 	It("Trains the model with the new data", func() {
					// 		orig := rm.Classifiers[0].TrainingData.Records

					// 		if bodyBytes != nil {
					// 			body = bytes.NewReader(bodyBytes)
					// 		} else {
					// 			body = nil
					// 		}
					// 		request = httptest.NewRequest(method, target, body)
					// 		request.Header.Add("Content-type", "application/json")
					// 		c = echo.New().NewContext(request, recorder)
					// 		handlers.TrainModelHandler(rm)(c)
					// 		Expect(rm.Classifiers[0].TrainingData.Records).NotTo(ConsistOf(orig))
					// 	})
					// })
				})
			})

			When("The body is not valid JSON", func() {
				BeforeEach(func() {
					bodyBytes, _ = os.ReadFile("../../fixtures/shorebirds_bad.json")
				})

				It("Returns a 400", func() {
					handlers.TrainModelHandler(rm)(c)
					resp := recorder.Result()
					Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				})
			})

			When("The body is valid JSON but has data inconsistency issues", func() {
				BeforeEach(func() {
					bodyBytes, _ = os.ReadFile("../../fixtures/shorebirds_badattributes.json")
				})

				It("Returns a 400", func() {
					handlers.TrainModelHandler(rm)(c)
					resp := recorder.Result()
					Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				})
			})
		})

		When("The model has not been set in the context", func() {
			It("Returns a 404", func() {
				handlers.TrainModelHandler(rm)(c)
				resp := recorder.Result()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
})
