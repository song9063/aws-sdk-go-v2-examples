package main

// https://aws.amazon.com/ko/getting-started/hands-on/design-a-database-for-a-mobile-app-with-dynamodb/5/
// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb@v1.2.1#Client.UpdateTable

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
	const indexName = "InvertedIndexEx"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(cfg)

	updateTableInput := &dynamodb.UpdateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			types.AttributeDefinition{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			types.AttributeDefinition{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{
			types.GlobalSecondaryIndexUpdate{
				Create: &types.CreateGlobalSecondaryIndexAction{
					IndexName: aws.String(indexName),
					KeySchema: []types.KeySchemaElement{
						types.KeySchemaElement{
							AttributeName: aws.String("SK"),
							KeyType:       types.KeyTypeHash,
						},
						types.KeySchemaElement{
							AttributeName: aws.String("PK"),
							KeyType:       types.KeyTypeRange,
						},
					},
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
					ProvisionedThroughput: &types.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(5),
						WriteCapacityUnits: aws.Int64(5),
					},
				},
			},
		},
	} // end of updateTableInput

	output, err := client.UpdateTable(context.TODO(), updateTableInput)
	if err != nil {
		log.Fatal(err)
	}
	objectDump(output)

}
