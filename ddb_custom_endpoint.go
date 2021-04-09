package main

// Custom endpoint(ex: DynamoDB local)
// https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/endpoints/
// https://aws.amazon.com/ko/getting-started/hands-on/design-a-database-for-a-mobile-app-with-dynamodb/4/
// https://docs.aws.amazon.com/ko_kr/cli/latest/reference/dynamodb/index.html
// https://www.aws.training/Details/eLearning?id=65040
// Exploring the DynamoDB API and the AWS SDKs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func loadJsonItems(itemFilePath string) ([]interface{}, error) {
	file, err := os.Open(itemFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var itemsJson []interface{} = nil
	err = json.Unmarshal(bytes, &itemsJson)
	if err != nil {
		return nil, err
	}
	return itemsJson, nil
}

func makeWriteRequestsFromJsonArray(jsonItems []interface{}) ([]types.WriteRequest, error) {
	arWriteRequests := make([]types.WriteRequest, 0, 25)
	for _, jsonItem := range jsonItems {
		itemForPut, err := attributevalue.MarshalMap(jsonItem)
		if err != nil {
			return nil, err
		}
		writeReq := types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: itemForPut,
			},
		}
		arWriteRequests = append(arWriteRequests, writeReq)
	}
	return arWriteRequests, nil
}
func BatchWrite(ddbClient *dynamodb.Client, strTableName string, arRequests []types.WriteRequest) (*dynamodb.BatchWriteItemOutput, error) {
	params := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			strTableName: arRequests,
		},
	}
	output, err := ddbClient.BatchWriteItem(context.TODO(), params)
	return output, err
}
func main() {
	const strTableName = "Movies"
	const strItemFileName = "moviedata.json"
	const maxItemsInBatch = 25

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path = path + string(os.PathSeparator) + strItemFileName
	fmt.Println(path)

	jsonItems, err := loadJsonItems(path)
	if err != nil {
		log.Fatal(err)
	}
	lengthOfItems := len(jsonItems)
	fmt.Printf("%d items are parsed.\n", lengthOfItems)

	// Connect to DB
	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				//PartitionID:   "aws",
				URL: "http://localhost:8000",
				//SigningRegion: "us-west-2",
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolver(customResolver))
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(cfg)

	// BatchWrite
	for startIndex := 0; startIndex < lengthOfItems; startIndex += maxItemsInBatch {
		endIndex := startIndex + maxItemsInBatch
		if endIndex > lengthOfItems {
			endIndex = lengthOfItems
		}

		requestArray, err := makeWriteRequestsFromJsonArray(jsonItems[startIndex:endIndex])
		output, err := BatchWrite(client, strTableName, requestArray)
		if err != nil {
			log.Fatal(err)
		}
		// /*for debug*/
		if output != nil {
			jsonOutput, _ := json.MarshalIndent(output, "", "\t")
			fmt.Println(string(jsonOutput))
		}
		time.Sleep(220 * time.Millisecond)
	}

}
