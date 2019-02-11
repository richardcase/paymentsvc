package service

/*
import (
	"net/http"
	"testing"

	httpdelivery "github.com/richardcase/paymentsvc/pkg/http"
	"github.com/richardcase/paymentsvc/pkg/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)



func TestCanFetchPayment(t *testing.T) {
	request := httpdelivery.Req{
		HTTPMethod:     "GET",
		PathParameters: map[string]string{"id": "11a7dc46-0176-4c91-a148-a578cc9fd3a2"},
	}
	h := &Handler{&MockPaymentRepository{}}
	router := httpdelivery.Router(h)
	response, err := router(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestCanCreatePayment(t *testing.T) {
	request := httpdelivery.Req{
		HTTPMethod: "POST",
		Body:       `{ "name": "test client", "description": "some test", "rate": 40 }`,
	}
	h := &Handler{&MockPaymentRepository{}}
	router := httpdelivery.Router(h)
	response, err := router(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)
}

func TestCanListPayments(t *testing.T) {
	request := httpdelivery.Req{
		HTTPMethod: "GET",
	}
	h := &Handler{&MockPaymentRepository{}}
	router := httpdelivery.Router(h)
	response, err := router(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestHandleInvalidJSON(t *testing.T) {
	request := httpdelivery.Req{
		HTTPMethod: "POST",
		Body:       "",
	}
	h := &Handler{&MockPaymentRepository{}}
	router := httpdelivery.Router(h)
	response, err := router(request)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

*/
