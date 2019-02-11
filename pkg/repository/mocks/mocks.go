package mocks

import _ "github.com/vektra/mockery"

//go:generate ${GOPATH}/bin/mockery -tags netgo -dir=../ -name=Repository -output=./
//go:generate ${GOPATH}/bin/mockery -tags netgo -dir=../../../vendor/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface -name=DynamoDBAPI -output=./
