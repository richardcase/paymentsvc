package repository

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/richardcase/paymentsvc/pkg/model"
)

// DynamoDbRepository is a payments repository implemented using DynamoDB
type DynamoDbRepository struct {
	table string
	conn  dynamodbiface.DynamoDBAPI
}

// NewDynamoDbRepository creates a new DynamoDbRepository
func NewDynamoDbRepository(tablename string, db dynamodbiface.DynamoDBAPI) (*DynamoDbRepository, error) {
	err := checkTable(db, tablename)
	if err != nil {
		return nil, err
	}

	return &DynamoDbRepository{
		conn:  db,
		table: tablename,
	}, nil
}

// GetAll will return all payments from dynamodb. Note: pagination will be added in the future.
func (r *DynamoDbRepository) GetAll() ([]*model.Payment, error) {
	var payments []*model.Payment

	results, err := r.conn.Scan(&dynamodb.ScanInput{
		TableName: aws.String(r.table),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "scanning table %s", r.table)
	}
	if err := dynamodbattribute.UnmarshalListOfMaps(results.Items, &payments); err != nil {
		return nil, errors.Wrap(err, "unmarshalling dynamo results")
	}

	return payments, nil
}

// GetByID will get a specific payment from DynamoDB
func (r *DynamoDbRepository) GetByID(id string) (*model.Payment, error) {
	var payment *model.Payment

	result, err := r.conn.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "getting payment from dynamo %s", id)
	}

	if result.Item == nil || len(result.Item) == 0 {
		return nil, nil
	}

	if err := dynamodbattribute.UnmarshalMap(result.Item, &payment); err != nil {
		return nil, errors.Wrap(err, "unmarshalling dynamo response to payment")
	}
	return payment, nil
}

// Create will create a new payment and save in dynamodb
func (r *DynamoDbRepository) Create(attributes *model.PaymentAttributes) (*model.Payment, error) {
	now := time.Now()
	payment := &model.Payment{
		ID:             uuid.New().String(),
		OrganisationID: uuid.New().String(),
		Version:        1,
		Attributes:     *attributes,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	av, err := dynamodbattribute.MarshalMap(payment)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling payment to synamo map")
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.table),
	}
	_, err = r.conn.PutItem(input)
	if err != nil {
		return nil, errors.Wrapf(err, "putting payment into table %s", r.table)
	}
	return payment, err
}

// Update will update an existing payment in DynamoDb
func (r *DynamoDbRepository) Update(id string, attributes *model.PaymentAttributes) (*model.Payment, error) {
	// Get current version
	payment, err := r.GetByID(id)
	if err != nil {
		return nil, errors.Wrapf(err, "getting existing payment %s", id)
	}

	av, err := dynamodbattribute.MarshalMap(attributes)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling payment attributes to synamo map")
	}

	newVersion := strconv.Itoa(payment.Version + 1)
	updatedAt := strconv.FormatInt(time.Now().Unix(), 10)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		UpdateExpression: aws.String("SET attributes = :attributes, version = :version, updated_at = :updated_at"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":attributes": {
				M: av,
			},
			":version": {
				N: aws.String(newVersion),
			},
			":updated_at": {
				N: aws.String(updatedAt),
			},
		},
		ReturnValues: aws.String(dynamodb.ReturnValueAllNew),
	}
	result, err := r.conn.UpdateItem(input)
	if err != nil {
		return nil, errors.Wrapf(err, "updating payment into table %s", r.table)
	}

	var updatedPayment *model.Payment
	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &updatedPayment); err != nil {
		return nil, errors.Wrap(err, "unmarshalling dynamo response to payment")
	}

	return updatedPayment, nil
}

// Delete will delete a specific payment from dynamodb
func (r *DynamoDbRepository) Delete(id string) error {

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	_, err := r.conn.DeleteItem(input)

	if err != nil {
		return errors.Wrapf(err, "deleting payment in dynamo %s", id)
	}

	return nil
}

func checkTable(conn dynamodbiface.DynamoDBAPI, tablename string) error {
	fmt.Printf("checking table %s exists", tablename)
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tablename),
	}
	result, err := conn.DescribeTable(input)
	if err != nil {
		return errors.Wrapf(err, "describing table %s", tablename)
	}

	// check the table is active
	if *result.Table.TableStatus != dynamodb.TableStatusActive {
		return fmt.Errorf("table %s status isn't ACTIVE but %s instead", tablename, *result.Table.TableStatus)
	}

	return nil
}
