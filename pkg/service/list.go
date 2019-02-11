package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/richardcase/paymentsvc/pkg/model"
)

// List is the handler for getting a list of payments. NOTE: no pagination at present
func (s *Service) List(c *gin.Context) {
	payments, err := s.repository.GetAll()
	if err != nil {
		errResp := &model.ErrorResponse{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	if payments == nil {
		paymentsEmpty := make([]model.Payment, 0)
		c.JSON(http.StatusOK, paymentsEmpty)
		return
	}

	c.JSON(http.StatusOK, payments)
}
