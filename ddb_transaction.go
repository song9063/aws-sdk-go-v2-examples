package main

// https://aws.amazon.com/ko/getting-started/hands-on/design-a-database-for-a-mobile-app-with-dynamodb/7/
// add_reaction.py

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
	const reactingUser = "john42"
	const reactingType = "sunglasses"
	const photoUser = "ppierce"
	const photoTimestamp = "2019-04-14T08:09:34"

	reaction := fmt.Sprintf("REACTION#%s#%s", reactingUser, reactingType)
	photo := fmt.Sprintf("PHOTO#%s#%s", photoUser, photoTimestamp)
	user := fmt.Sprintf("USER#%s", photoUser)

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
						"PK":           &types.AttributeValueMemberS{Value: reaction},
						"SK":           &types.AttributeValueMemberS{Value: photo},
						"reactingUser": &types.AttributeValueMemberS{Value: reactingUser},
						"reactionType": &types.AttributeValueMemberS{Value: reactingType},
						"photo":        &types.AttributeValueMemberS{Value: photo},
						"timestamp":    &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
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
						"SK": &types.AttributeValueMemberS{Value: photo},
					},
					UpdateExpression: aws.String("SET reactions.#t = reactions.#t + :i"),
					ExpressionAttributeNames: map[string]string{
						"#t": reactingType,
					},
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
