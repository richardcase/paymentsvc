package service_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/richardcase/paymentsvc/pkg/repository/mocks"
	"github.com/richardcase/paymentsvc/pkg/service"
)

var _ = Describe("Create Payment", func() {
	var (
		r         *mocks.Repository
		router    *gin.Engine
		svc       *service.Service
		err       error
		bodyBytes []byte
		recorder  *httptest.ResponseRecorder
		res       *http.Response
		req       *http.Request
	)

	Describe("when calling POST /payments", func() {

		Describe("with a valid request body", func() {
			BeforeEach(func() {
				bodyBytes, _ = ioutil.ReadFile("testdata/create_body_req.golden")

				r = &mocks.Repository{}
				payment := createPayment()
				attr := createPaymentAttributes()
				r.On("Create", attr).Return(payment, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("POST", "/payments", strings.NewReader(string(bodyBytes)))
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(recorder, req)
				res = recorder.Result()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have HTTP status code 201 (Created)", func() {
				Expect(res.StatusCode).Should(Equal(http.StatusCreated))
			})
			It("should have a Location header", func() {
				hdr := res.Header["Location"]
				Expect(hdr).To(Not(BeNil()))
				Expect(len(hdr)).To(Equal(1))

				loc := hdr[0]
				Expect(strings.HasPrefix(loc, "/")).To(BeTrue())
			})

			It("should have created the payment", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "Create", 1)).To(BeTrue())
			})

		})

		Describe("with an invalid request body", func() {
			BeforeEach(func() {
				bodyBytes, _ = ioutil.ReadFile("testdata/create_invalidbody_req.golden")

				r = &mocks.Repository{}
				payment := createPayment()
				attr := createPaymentAttributes()
				r.On("Create", attr).Return(payment, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("POST", "/payments", strings.NewReader(string(bodyBytes)))
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
			It("should not have stored the payment", func() {
				Expect(r.AssertNotCalled(GinkgoT(), "Create")).To(BeTrue())
			})

		})

		Describe("with a valid request body and datastore errors", func() {
			BeforeEach(func() {
				bodyBytes, _ = ioutil.ReadFile("testdata/create_body_req.golden")

				r = &mocks.Repository{}
				attr := createPaymentAttributes()
				r.On("Create", attr).Return(nil, errors.New("error saving payment"))

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("POST", "/payments", strings.NewReader(string(bodyBytes)))
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(recorder, req)
				res = recorder.Result()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have HTTP status code 500 (Internal Server Error)", func() {
				Expect(res.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})

	})
})
