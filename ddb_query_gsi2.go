package main

//https://aws.amazon.com/ko/getting-started/hands-on/design-a-database-for-a-mobile-app-with-dynamodb/5/
// find_following_for_user.py

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func objectDump(obj interface{}) {
	jsonOutput, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(jsonOutput))
}

func main() {
	const tableName = "quick-photos"
	const indexName = "InvertedIndex"
	const userName = "haroldwatkins"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(cfg)

	qInput := dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String(indexName),
		KeyConditionExpression: aws.String("SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("#FRIEND#%s", userName),
			},
		},
		ScanIndexForward: aws.Bool(true),
	}
	objectDump(qInput)
	output, err := client.Query(context.TODO(), &qInput)

	objectDump(output.Items)
}
