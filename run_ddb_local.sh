#!/bin/bash
# Execute DynamoDB Local
# aws dynamodb list-tables --endpoint-url http://localhost:8000
java -Djava.library.path=./dynamodb_local_latest/DynamoDBLocal_lib -jar ./dynamodb_local_latest/DynamoDBLocal.jar -sharedDb
