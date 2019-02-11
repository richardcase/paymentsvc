package service_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/richardcase/paymentsvc/pkg/model"
	"github.com/richardcase/paymentsvc/pkg/repository/mocks"
	"github.com/richardcase/paymentsvc/pkg/service"
)

var _ = Describe("Update Payment", func() {
	var (
		r         *mocks.Repository
		router    *gin.Engine
		svc       *service.Service
		err       error
		paymentId string
		bodyBytes []byte
		recorder  *httptest.ResponseRecorder
		res       *http.Response
		req       *http.Request
	)

	BeforeEach(func() {
		paymentId = "cae9aa62-0ea1-432b-baee-c0ff4b1d889e"
		bodyBytes, _ = ioutil.ReadFile("testdata/update_body_req.golden")
	})

	Describe("when calling PUT /payments/{id}", func() {

		Describe("with an incorrectly formatted id", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				currentVerion := createPayment()
				newVersion := createPaymentWithOverrides(500.00, 2)
				updatedAttribs := createPaymentAttributesWithOveride(500.00)

				r.On("GetByID", paymentId).Return(currentVerion, nil)
				r.On("Update", paymentId, updatedAttribs).Return(newVersion, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("PUT", "/payments/ABCD", strings.NewReader(string(bodyBytes)))
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
			It("should not have updated the payment", func() {
				Expect(r.AssertNotCalled(GinkgoT(), "Update")).To(BeTrue())
			})
			It("should not have fetched existing payment", func() {
				Expect(r.AssertNotCalled(GinkgoT(), "GetByID")).To(BeTrue())
			})
		})

		Describe("with an id for a payment that doesn't exist", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				newVersion := createPaymentWithOverrides(500.00, 2)
				updatedAttribs := createPaymentAttributesWithOveride(500.00)

				r.On("GetByID", paymentId).Return(nil, nil)
				r.On("Update", paymentId, updatedAttribs).Return(newVersion, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("PUT", "/payments/"+paymentId, strings.NewReader(string(bodyBytes)))
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
			It("should not have updated the payment", func() {
				Expect(r.AssertNotCalled(GinkgoT(), "Update")).To(BeTrue())
			})
			It("should have tried to get the payment", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetByID", 1)).To(BeTrue())
			})
		})

		Describe("with an id for a payment that exists and valid body", func() {
			var (
				currentVersion *model.Payment
				newVersion     *model.Payment
			)
			BeforeEach(func() {
				r = &mocks.Repository{}
				currentVersion = createPayment()
				newVersion = createPaymentWithOverrides(500.00, 2)
				updatedAttribs := createPaymentAttributesWithOveride(500.00)

				r.On("GetByID", paymentId).Return(currentVersion, nil)
				r.On("Update", paymentId, updatedAttribs).Return(newVersion, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("PUT", "/payments/"+paymentId, strings.NewReader(string(bodyBytes)))
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
			It("should have updated the payment", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "Update", 1)).To(BeTrue())
			})
			It("should have got the payment", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetByID", 1)).To(BeTrue())
			})
			It("should have returned JSON body with 1 payment", func() {
				g, err := ioutil.ReadFile("testdata/update_body_resp.golden")
				if err != nil {
					Fail("error reading golden file")
				}
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				Expect(string(bodyBytes)).To(MatchJSON(string(g)))
			})
		})

		Describe("with an id for a payment that exists and invalid body", func() {
			var (
				currentVersion *model.Payment
				newVersion     *model.Payment
			)
			BeforeEach(func() {
				bodyBytes, _ = ioutil.ReadFile("testdata/update_invalidbody_req.golden")

				r = &mocks.Repository{}
				currentVersion = createPayment()
				newVersion = createPaymentWithOverrides(500.00, 2)
				updatedAttribs := createPaymentAttributesWithOveride(0.00)

				r.On("GetByID", paymentId).Return(currentVersion, nil)
				r.On("Update", paymentId, updatedAttribs).Return(newVersion, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("PUT", "/payments/"+paymentId, strings.NewReader(string(bodyBytes)))
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
			It("should not have updated the payment", func() {
				Expect(r.AssertNotCalled(GinkgoT(), "Update")).To(BeTrue())
			})
		})
	})
})
