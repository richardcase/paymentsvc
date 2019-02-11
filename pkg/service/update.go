package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/richardcase/paymentsvc/pkg/model"
)

// Update is the handler for updating a payment
func (s *Service) Update(c *gin.Context) {
	paymentID := c.Param("id")
	id, err := uuid.Parse(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid payment id",
			"error":   err.Error(),
		})
		return
	}

	payment, err := s.repository.GetByID(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error querying for payment",
			"error":   err.Error(),
		})
		return
	}
	if payment == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Couldn't find payment %s", id.String()),
			"error":   "NotFound",
		})
		return
	}

	var attributes model.PaymentAttributes
	err = c.ShouldBindJSON(&attributes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed",
			"error":   err.Error(),
		})
		return
	}

	payment, err = s.repository.Update(paymentID, &attributes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error updating payment in datastore",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, payment)
}
