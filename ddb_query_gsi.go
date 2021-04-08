package main

// https://aws.amazon.com/ko/getting-started/hands-on/design-a-database-for-a-mobile-app-with-dynamodb/5/
// fetch_photo_and_reactions.py

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const tableName = "quick-photos"
const indexName = "InvertedIndex"

func objectDump(obj interface{}) {
	jsonOutput, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(jsonOutput))
}

func reverse(ar *[]map[string]types.AttributeValue) {
	arr := *ar
	for i := len(arr)/2 - 1; i >= 0; i-- {
		opp := len(arr) - 1 - i
		arr[i], arr[opp] = arr[opp], arr[i]
	}
}

func fetchPhotoAndReactions(pClient *dynamodb.Client, userName, timeStamp string) (*dynamodb.QueryOutput, error) {
	if pClient == nil {
		return nil, errors.New("pClient can't be nil")
	}
	if len(userName) < 1 {
		return nil, errors.New("Required: userName")
	}
	if len(timeStamp) < 1 {
		return nil, errors.New("Required: timeStamp")
	}

	qInput := dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String(indexName),
		KeyConditionExpression: aws.String("SK = :sk AND PK BETWEEN :reactions AND :user"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("PHOTO#%s#%s", userName, timeStamp),
			},
			":user": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("USER$"),
			},
			":reactions": &types.AttributeValueMemberS{
				Value: "REACTION#",
			},
		},
		ScanIndexForward: aws.Bool(true),
	}
	objectDump(qInput)
	output, err := pClient.Query(context.TODO(), &qInput)
	return output, err
}

func main() {
	userName := "david25"
	timeStamp := "2019-03-02T09:11:30"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	qOut, err := fetchPhotoAndReactions(client, userName, timeStamp)
	if err != nil {
		log.Fatal(err)
	}

	var photos []map[string]types.AttributeValue = nil
	if qOut.Count > 0 {
		photos = qOut.Items
	}
	reverse(&photos)
	objectDump(photos)
}
