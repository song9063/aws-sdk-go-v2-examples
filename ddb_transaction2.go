package main

// https://aws.amazon.com/ko/getting-started/hands-on/design-a-database-for-a-mobile-app-with-dynamodb/7/
// follow_user.py

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	const followedUser = "tmartinez"
	const followingUser = "john42"

	user := fmt.Sprintf("USER#%s", followedUser)
	friend := fmt.Sprintf("#FRIEND#%s", followingUser)
	userMetadata := fmt.Sprintf("#METADATA#%s", followedUser)
	friendUser := fmt.Sprintf("USER#%s", followingUser)
	friendMetadata := fmt.Sprintf("#METADATA#%s", followingUser)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(cfg)

	tOut, err := client.TransactWriteItems(context.TODO(), &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			types.TransactWriteItem{
				Put: &types.Put{
					TableName: aws.String(tableName),
					Item: map[string]types.AttributeValue{
						"PK":            &types.AttributeValueMemberS{Value: user},
						"SK":            &types.AttributeValueMemberS{Value: friend},
						"followedUser":  &types.AttributeValueMemberS{Value: followedUser},
						"followingUser": &types.AttributeValueMemberS{Value: followingUser},
						"timestamp":     &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
					},
					ConditionExpression:                 aws.String("attribute_not_exists(SK)"),
					ReturnValuesOnConditionCheckFailure: types.ReturnValuesOnConditionCheckFailureAllOld,
				},
			},
			types.TransactWriteItem{
				Update: &types.Update{
					TableName: aws.String(tableName),
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: user},
						"SK": &types.AttributeValueMemberS{Value: userMetadata},
					},
					UpdateExpression: aws.String("SET followers = followers + :i"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":i": &types.AttributeValueMemberN{Value: "1"},
					},
					ReturnValuesOnConditionCheckFailure: types.ReturnValuesOnConditionCheckFailureAllOld,
				},
			},
			types.TransactWriteItem{
				Update: &types.Update{
					TableName: aws.String(tableName),
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: friendUser},
						"SK": &types.AttributeValueMemberS{Value: friendMetadata},
					},
					UpdateExpression: aws.String("SET following = following + :i"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":i": &types.AttributeValueMemberN{Value: "1"},
					},
					ReturnValuesOnConditionCheckFailure: types.ReturnValuesOnConditionCheckFailureAllOld,
				},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
	objectDump(tOut)
}
