package service_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/richardcase/paymentsvc/pkg/model"
)

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

func createPaymentAttributes() *model.PaymentAttributes {
	payment := createPayment()

	return &payment.Attributes
}

func createPaymentAttributesWithOveride(amount float32) *model.PaymentAttributes {
	payment := createPaymentWithOverrides(amount, 1)

	return &payment.Attributes
}

func createPaymentWithOverrides(amount float32, version int) *model.Payment {
	responsePayment := createPayment()
	responsePayment.Attributes.Amount = amount
	responsePayment.Version = version

	return responsePayment
}

func createPayment() *model.Payment {
	responseBytes := getPaymentResponseTemplate()
	responsePayment := &model.Payment{}
	err := json.Unmarshal(responseBytes, responsePayment)
	if err != nil {
		Fail(fmt.Sprintf("error creating payment: %s", err.Error()))
	}

	return responsePayment
}

func getPaymentResponseTemplate() []byte {
	responseBytes, err := ioutil.ReadFile("testdata/create_payment_template.golden")
	if err != nil {
		Fail(fmt.Sprintf("error reading golden file: %s", err.Error()))
	}
	return responseBytes
}
