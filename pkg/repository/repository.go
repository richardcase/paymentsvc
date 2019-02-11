package repository

import (
	"github.com/richardcase/paymentsvc/pkg/model"
)

// Repository represents a store of payments
type Repository interface {
	GetAll() ([]*model.Payment, error) //TODO: Support pagination
	GetByID(id string) (*model.Payment, error)
	Create(attributes *model.PaymentAttributes) (*model.Payment, error)
	Update(id string, attributes *model.PaymentAttributes) (*model.Payment, error)
	Delete(id string) error
}
