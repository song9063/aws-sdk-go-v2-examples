package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func objectDump(obj interface{}) {
	jsonOutput, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(jsonOutput))
}
func main() {
	const tableName = "quick-photos"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(cfg)

	dOut, err := client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		log.Fatal(err)
	}
	objectDump(dOut)
}
