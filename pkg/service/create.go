package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/richardcase/paymentsvc/pkg/model"
)

// Create is the handler for payment creation
func (s *Service) Create(c *gin.Context) {
	var attributes model.PaymentAttributes
	err := c.ShouldBindJSON(&attributes)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed",
			"error":   err.Error(),
		})
		return
	}

	payment, err := s.repository.Create(&attributes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error saving payment to datastore",
			"error":   err.Error(),
		})
		return
	}

	location := fmt.Sprintf("/payments/%s", payment.ID)

	c.Header("Location", location)
	c.JSON(http.StatusCreated, payment)
}
