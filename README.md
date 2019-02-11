

## Running locally

docker run --rm -p 8000:8000 --name dynamodb amazon/dynamodb-local -jar DynamoDBLocal.jar -inMemory -sharedDb

1st window
AWS_PROFILE=mon-dev make start-local-db

2nd window
AWS_PROFILE=mon-dev make populate-local-db
AWS_PROFILE=mon-dev make start-local-func 

3rd window
make integration-test-local

## Future

docker network create lambda-local
docker run --rm -p 8000:8000 --network lambda-local --name dynamodb amazon/dynamodb-local -jar DynamoDBLocal.jar -inMemory -sharedDb

docker run --network lambda-local -v $(pwd)/bin:/var/task -i -t --rm --env-file deployment/env/integration.env  lambci/lambda:go1.x payments

cat test/integration/full_lifecycle_test.go | docker run --network lambda-local -v $(pwd)/bin:/var/task -i -t --rm --env-file deployment/env/integration.env -e DOCKER_LAMBDA_USE_STDIN=1  lambci/lambda:go1.x payments


sam local start-api --region eu-west-2 --profile mon-dev -l out.log --env-vars env/local_mac.json

