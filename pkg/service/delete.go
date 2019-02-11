package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Delete is the handler for payment deletion
func (s *Service) Delete(c *gin.Context) {
	paymentID := c.Param("id")
	id, err := uuid.Parse(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid payment id",
			"error":   err.Error(),
		})
		return
	}

	// Check the payment exists first
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

	err = s.repository.Delete(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error deleteing for payment",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
