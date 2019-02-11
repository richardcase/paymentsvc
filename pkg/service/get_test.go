package service_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/richardcase/paymentsvc/pkg/model"
	"github.com/richardcase/paymentsvc/pkg/repository/mocks"
	"github.com/richardcase/paymentsvc/pkg/service"
)

var _ = Describe("Get Payment", func() {
	var (
		r         *mocks.Repository
		router    *gin.Engine
		svc       *service.Service
		err       error
		paymentId string
		recorder  *httptest.ResponseRecorder
		res       *http.Response
		req       *http.Request
	)

	BeforeEach(func() {
		paymentId = "cae9aa62-0ea1-432b-baee-c0ff4b1d889e"
	})

	Describe("when calling GET /payments/{id}", func() {

		Describe("with an incorrectly formatted id", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("GetByID", paymentId).Return(nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("GET", "/payments/ABCD", nil)
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(recorder, req)
				res = recorder.Result()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have HTTP status code 400 (BadRequest)", func() {
				Expect(res.StatusCode).Should(Equal(http.StatusBadRequest))
			})
			It("should not have deleted the payment", func() {
				Expect(r.AssertNotCalled(GinkgoT(), "GetByID")).To(BeTrue())
			})
		})

		Describe("with for a payment that doesn't exist", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("GetByID", paymentId).Return(nil, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("GET", "/payments/"+paymentId, nil)
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(recorder, req)
				res = recorder.Result()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have HTTP status code 404 (NotFound)", func() {
				Expect(res.StatusCode).Should(Equal(http.StatusNotFound))
			})
			It("should have tried to get the payment", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetByID", 1)).To(BeTrue())
			})
		})

		Describe("with an id for a payment that exists", func() {
			var (
				payment *model.Payment
			)
			BeforeEach(func() {
				r = &mocks.Repository{}
				payment = createPayment()
				r.On("GetByID", paymentId).Return(payment, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("GET", "/payments/"+paymentId, nil)
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(recorder, req)
				res = recorder.Result()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have HTTP status code 200 (OK)", func() {
				Expect(res.StatusCode).Should(Equal(http.StatusOK))
			})
			It("should have got the payment", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetByID", 1)).To(BeTrue())
			})
			It("should have a JSON body returned with 1 payemnt", func() {
				g := getPaymentResponseTemplate()
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				Expect(string(bodyBytes)).To(MatchJSON(string(g)))
			})
		})

		Describe("with an id for a payment that exists BUT the datastore errors", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("GetByID", paymentId).Return(nil, errors.New("datastore had an error"))

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("GET", "/payments/"+paymentId, nil)
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(recorder, req)
				res = recorder.Result()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have HTTP status code 500 (InternalServerError)", func() {
				Expect(res.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})

	})
})
