package main

import (
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	dbRegion   string
	dbEndpoint string
	tableName  string
)

func init() {
	flag.StringVar(&dbRegion, "region", "eu-west-2", "the AWS region")
	flag.StringVar(&dbEndpoint, "endpoint", "http://127.0.0.1:8000", "the url of the local dynamodb")
	flag.StringVar(&tableName, "table", "payments", "the name of the table")
}

func main() {
	log.Println("starting dynamodb populate util")
	flag.Parse()

	if dbRegion == "" {
		log.Panic("You must supply a region")
		return
	}
	if dbEndpoint == "" {
		log.Panic("You must supply a db endpoint")
		return
	}
	if tableName == "" {
		log.Panic("You must supply a table name")
		return
	}

	config := &aws.Config{
		Region:   aws.String(dbRegion),
		Endpoint: aws.String(dbEndpoint),
	}

	sess, err := session.NewSession(config)
	if err != nil {
		log.Panicf("Erroring creating AWS session %s", err.Error())
	}

	db := dynamodb.New(sess)

	params := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			&dynamodb.KeySchemaElement{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}

	createOut, err := db.CreateTable(params)
	if err != nil {
		log.Panicf("Error creating table %s", err.Error())
		return
	}
	log.Printf("Created table")
	log.Printf("%v", createOut)
}
