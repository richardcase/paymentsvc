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

var _ = Describe("Delete Payment", func() {
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

	Describe("when calling POST /payments", func() {

		Describe("with an incorrectly formatter id", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("Delete", paymentId).Return(nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("DELETE", "/payments/ABCD", nil)
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
				Expect(r.AssertNotCalled(GinkgoT(), "Delete")).To(BeTrue())
			})
		})

		Describe("with for a payment that doesn't exist", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("Delete", paymentId).Return(nil)
				r.On("GetByID", paymentId).Return(nil, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("DELETE", "/payments/"+paymentId, nil)
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
			It("should not have deleted the payment", func() {
				Expect(r.AssertNotCalled(GinkgoT(), "Delete")).To(BeTrue())
			})
			It("should have tried to get the payment", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetByID", 1)).To(BeTrue())
			})
		})

		Describe("with for a payment that exists", func() {
			var (
				payment *model.Payment
			)
			BeforeEach(func() {
				r = &mocks.Repository{}
				payment = createPayment()
				r.On("Delete", paymentId).Return(nil)
				r.On("GetByID", paymentId).Return(payment, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("DELETE", "/payments/"+paymentId, nil)
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(recorder, req)
				res = recorder.Result()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have HTTP status code 204 (NoContent)", func() {
				Expect(res.StatusCode).Should(Equal(http.StatusNoContent))
			})
			It("should have deleted the payment", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "Delete", 1)).To(BeTrue())
			})
			It("should have got the payment", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetByID", 1)).To(BeTrue())
			})
			It("should have not body returned", func() {
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				Expect(string(bodyBytes)).To(BeEmpty())
			})
		})

		Describe("with a payment that exists and datastore errors", func() {
			var (
				payment *model.Payment
			)
			BeforeEach(func() {
				r = &mocks.Repository{}
				payment = createPayment()
				r.On("Delete", paymentId).Return(errors.New("error deleting payment"))
				r.On("GetByID", paymentId).Return(payment, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("DELETE", "/payments/"+paymentId, nil)
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
