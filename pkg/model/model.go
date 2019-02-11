package model

import "time"

// Payment represents a payment
type Payment struct {
	ID             string            `json:"id" binding:"required,uuid" dynamodbav:"id,string"`
	OrganisationID string            `json:"organisation_id" binding:"required,uuid" dynamodbav:"organisation_id,string"`
	Version        int               `json:"version" binding:"required" dynamodbav:"version"`
	Attributes     PaymentAttributes `json:"attributes" binding:"required" dynamodbav:"attributes"`
	CreatedAt      time.Time         `json:"created_at"  dynamodbav:"created_at,unixtime"`
	UpdatedAt      time.Time         `json:"updated_at"  dynamodbav:"updated_at,unixtime"`
}

// PaymentAttributes represents the attributes of a payment
type PaymentAttributes struct {
	Amount           float32 `json:"amount" binding:"required,gt=0"  dynamodbav:"amount,float"`
	BeneficiaryParty Party   `json:"beneficiary_party" binding:"required" dynamodbav:"beneficiary_party"`
	Currency         string  `json:"currency" binding:"required,len=3" dynamodbav:"currency,string"`
	DebtorParty      Party   `json:"debtor_party" binding:"required" dynamodbav:"debtor_party"`
	PaymentScheme    string  `json:"payment_scheme" binding:"required" dynamodbav:"payment_scheme"`
	PaymentType      string  `json:"payment_type" binding:"required" dynamodbav:"payment_type"`
}

// Party represents a payment party
type Party struct {
	AccountName   string `json:"account_name,omitempty" binding:"required" dynamodbav:"account_name"`
	AccountNumber string `json:"account_number,omitempty" binding:"required" dynamodbav:"account_number"`
}
