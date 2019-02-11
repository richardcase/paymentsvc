package service

import (
	"github.com/gin-gonic/gin"
	"github.com/richardcase/paymentsvc/pkg/repository"
)

// Service represent the payments service
type Service struct {
	repository repository.Repository
}

// New creates a new instance of the payments Service
func New(repo repository.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

// SetupRouter specifies all the route/handler mappings
func SetupRouter(s *Service) *gin.Engine {
	r := gin.Default()
	r.GET("/payments", s.List)
	r.GET("/payments/:id", s.Get)
	r.POST("/payments", s.Create)
	r.PUT("/payments/:id", s.Update)
	r.DELETE("/payments/:id", s.Delete)

	gin.SetMode(gin.ReleaseMode)

	return r
}
