// +build integration

package integration_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

var _ = Describe("(Integration) Full payment lifecycle test", func() {
	var (
		baseUrl   string
		res       *http.Response
		req       *http.Request
		err       error
		paymentId string
	)

	BeforeSuite(func() {
		if runLocal {
			baseUrl = "http://127.0.0.1:9000"
		} else {
			Fail("Running against AWS not implemented yet")
		}
	})

	Describe("when we have a payments service", func() {
		//TODO: start payments service using gexec

		Describe("when calling GET /payments", func() {
			It("it should execute successfully", func() {
				url := fmt.Sprintf("%s/payments", baseUrl)
				res, err = http.Get(url)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				bodyBytes, _ := ioutil.ReadAll(res.Body)
				Expect(bodyBytes).Should(MatchJSON(getGoldenFile("list_empty_resp")))
			})

			Describe("when calling POST /payments with a valid body", func() {
				It("should execute succesfully", func() {
					url := fmt.Sprintf("%s/payments", baseUrl)
					client := &http.Client{}
					req, err = http.NewRequest("POST", url, bytes.NewBuffer(getGoldenFile("create_body_req")))
					res, err = client.Do(req)

					Expect(err).NotTo(HaveOccurred())
					Expect(res.StatusCode).To(Equal(http.StatusCreated))

					hdr := res.Header["Location"]
					Expect(hdr).To(Not(BeNil()))
					Expect(len(hdr)).To(Equal(1))

					loc := hdr[0]
					locParts := strings.Split(loc, "/")
					Expect(locParts[1]).To(Equal("payments"))
					paymentId, err = uuid.Parse(locParts[2])
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("when calling GET /payments/{id} with the created id", func() {
					It("should execute sucessfully", func() {
						url := fmt.Sprintf("%s/payments/%s", baseUrl, paymentId.String())
						res, err = http.Get(url)

						Expect(err).NotTo(HaveOccurred())
						Expect(res.StatusCode).To(Equal(http.StatusOK))
						// Only check a few of the returned parameters
						bodyBytes, _ := ioutil.ReadAll(res.Body)
						version := gjson.Get(string(bodyBytes), "version")
						Expect(version.String()).To(Equal("1"))
						//TODO: add more checks to the data
					})

					Describe("when calling PUT /payments/{id} with a valid body and payment id", func() {
						It("should execute succesfullly and update the payment", func() {
							url := fmt.Sprintf("%s/payments/%s", baseUrl, paymentId.String())
							client := &http.Client{}
							req, err = http.NewRequest("PUT", url, bytes.NewBuffer(getGoldenFile("update_body_req")))
							res, err = client.Do(req)

							Expect(err).NotTo(HaveOccurred())
							Expect(res.StatusCode).To(Equal(http.StatusOK))

							bodyBytes, _ := ioutil.ReadAll(res.Body)
							version := gjson.Get(string(bodyBytes), "version")
							Expect(version.String()).To(Equal("2"))
							amount := gjson.Get(string(bodyBytes), "attributes.amount")
							Expect(amount.String()).To(Equal("500"))

						})

						Describe("when calling GET /payments", func() {
							It("should execute succesfully and return 1 payment", func() {
								url := fmt.Sprintf("%s/payments", baseUrl)
								res, err = http.Get(url)

								Expect(err).NotTo(HaveOccurred())
								Expect(res.StatusCode).To(Equal(http.StatusOK))
								bodyBytes, _ := ioutil.ReadAll(res.Body)
								body := gjson.Parse(string(bodyBytes))
								Expect(body.IsArray()).To(BeTrue())
								Expect(len(body.Array())).To(Equal(1))
								retId := body.Array()[0].Get("id")
								Expect(retId.String()).To(Equal(paymentId.String()))
								version := body.Array()[0].Get("version")
								Expect(version.String()).To(Equal("2"))
								amount := body.Array()[0].Get("attributes.amount")
								Expect(amount.String()).To(Equal("500"))
							})

							Describe("when calling DELETE /payments/{id} with a valid payment id", func() {
								It("should execute sucessfully", func() {
									url := fmt.Sprintf("%s/payments/%s", baseUrl, paymentId.String())
									client := &http.Client{}
									req, err = http.NewRequest("DELETE", url, nil)
									res, err = client.Do(req)

									Expect(err).NotTo(HaveOccurred())
									Expect(res.StatusCode).To(Equal(http.StatusNoContent))
								})

								Describe("when calling GET /payments/{id} with the id of the deleted payment", func() {
									It("should execute succesfull and return a 404", func() {
										url := fmt.Sprintf("%s/payments/%s", baseUrl, paymentId.String())
										res, err = http.Get(url)

										Expect(err).NotTo(HaveOccurred())
										Expect(res.StatusCode).To(Equal(http.StatusNotFound))
									})
								})

								Describe("when calling GET /payments", func() {
									It("it should execute successfully", func() {
										url := fmt.Sprintf("%s/payments", baseUrl)
										res, err = http.Get(url)

										Expect(err).NotTo(HaveOccurred())
										Expect(res.StatusCode).To(Equal(http.StatusOK))
										bodyBytes, _ := ioutil.ReadAll(res.Body)
										Expect(bodyBytes).Should(MatchJSON(getGoldenFile("list_empty_resp")))
									})
								})
							})
						})
					})
				})
			})
		})

	})
})

func getGoldenFile(name string) []byte {
	path := fmt.Sprintf("testdata/%s.golden", name)
	g, err := ioutil.ReadFile(path)
	if err != nil {
		Fail(fmt.Sprintf("failed reading golden file: %s", path))
	}
	return g
}
