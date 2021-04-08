package main

// https://aws.amazon.com/ko/getting-started/hands-on/design-a-database-for-a-mobile-app-with-dynamodb/4/

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

func objectDump(obj interface{}) {
	jsonOutput, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(jsonOutput))
}

func fetchUserAndPhotos(pClient *dynamodb.Client, userName string) (*dynamodb.QueryOutput, error) {
	if pClient == nil {
		return nil, errors.New("pClient can't be nil")
	}
	if len(userName) < 1 {
		return nil, errors.New("Required: userName")
	}
	qInput := dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("PK=:pk AND SK BETWEEN :metadata AND :photos"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("USER#%s", userName),
			},
			":metadata": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("#METADATA#%s", userName),
			},
			":photos": &types.AttributeValueMemberS{
				Value: "PHOTO$",
			},
		},
		ScanIndexForward: aws.Bool(true),
	}
	output, err := pClient.Query(context.TODO(), &qInput)
	return output, err
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	qOut, err := fetchUserAndPhotos(client, "john42")
	if err != nil {
		log.Fatal(err)
	}

	var user interface{} = nil
	var photos []map[string]types.AttributeValue = nil
	if qOut.Count > 0 {
		user = qOut.Items[0]
	}

	if qOut.Count > 1 {
		photos = qOut.Items[1:]
	}
	objectDump(user)
	objectDump(photos)
}
