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

var _ = Describe("List Payment", func() {
	var (
		r        *mocks.Repository
		router   *gin.Engine
		svc      *service.Service
		err      error
		recorder *httptest.ResponseRecorder
		res      *http.Response
		req      *http.Request
	)

	Describe("when calling GET /payments", func() {

		Describe("with no existing payments", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				payments := make([]*model.Payment, 0)
				r.On("GetAll").Return(payments, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("GET", "/payments", nil)
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
			It("should have called GetAll on the repo", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetAll", 1)).To(BeTrue())
			})
			It("should have an empty JSON array as body", func() {
				g, err := ioutil.ReadFile("testdata/list_empty_resp.golden")
				if err != nil {
					GinkgoT().Fatalf("failed reading .golden: %s", err)
				}
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				Expect(string(bodyBytes)).Should(MatchJSON(string(g)))
			})
		})

		Describe("with an existing payment", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				payments := make([]*model.Payment, 1)
				payments[0] = createPayment()
				r.On("GetAll").Return(payments, nil)

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("GET", "/payments", nil)
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
			It("should have called GetAll on the repo", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetAll", 1)).To(BeTrue())
			})
			It("should have an JSON array with 1 payment as body", func() {
				g, err := ioutil.ReadFile("testdata/list_single_resp.golden")
				if err != nil {
					GinkgoT().Fatalf("failed reading .golden: %s", err)
				}
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				Expect(string(bodyBytes)).Should(MatchJSON(string(g)))
			})
		})

		Describe("with an existing payment BUT datastore errors", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("GetAll").Return(nil, errors.New("datastore error"))

				svc = service.New(r)
				router = service.SetupRouter(svc)
				recorder = httptest.NewRecorder()

				req, err = http.NewRequest("GET", "/payments", nil)
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
			It("should have called GetAll on the repo", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetAll", 1)).To(BeTrue())
			})
		})
	})
})
