package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/pkg/errors"

	"github.com/richardcase/paymentsvc/pkg/repository"
	"github.com/richardcase/paymentsvc/pkg/service"
)

var (
	ginLambda *ginadapter.GinLambda
	s         *service.Service
	version   = "dev"
	commit    = "none"
	date      = "unknown"
)

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		log.Printf("Gin cold start")
		r := service.SetupRouter(s)
		ginLambda = ginadapter.New(r)
	}

	return ginLambda.Proxy(req)
}

func main() {
	fmt.Println("starting payments service")
	fmt.Printf("Version: %s", version)
	fmt.Printf("Build Date: %s", date)
	fmt.Printf("Git Commit: %s", commit)

	region := os.Getenv("REGION")
	table := os.Getenv("DB_TABLE")
	dbEndpoint := os.Getenv("DB_ENDPOINT_OVERRIDE")
	runLocal := os.Getenv("RUN_LOCAL")

	fmt.Printf("\tRegion=%s\n", region)
	fmt.Printf("\tTable=%s\n", table)
	fmt.Printf("\tDbEndpoint=%s\n", dbEndpoint)
	fmt.Printf("\tRunLocal=%s\n", runLocal)

	// Create connection to dynamo
	db, err := createDynamoDbConn(region, dbEndpoint)
	if err != nil {
		log.Panicf("error connecting to dynamodb: %s", err.Error())
	}

	// Create the repo
	repo, err := repository.NewDynamoDbRepository(table, db)
	if err != nil {
		log.Panic(err)
	}

	// Create the service instance, with the repository
	s = service.New(repo)

	// Start the Lambda process
	if runLocal == "" || runLocal != "true" {
		lambda.Start(Handler)
	} else {
		r := service.SetupRouter(s)
		err := r.Run(":9000")
		if err != nil {
			log.Panicf("error running functions: %s", err.Error())
		}
	}
}

func createDynamoDbConn(region string, endpoint string) (*dynamodb.DynamoDB, error) {
	config := &aws.Config{
		Region: aws.String(region),
	}
	if endpoint != "" {
		config.Endpoint = aws.String(endpoint)
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, errors.Wrapf(err, "creating new session in region %s", region)
	}
	db, err := dynamodb.New(sess), nil
	if err != nil {
		return nil, errors.Wrapf(err, "creating dynamodb instance using endpoint %s", endpoint)
	}
	return db, nil
}
