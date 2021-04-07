package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func makeWriteRequests(itemFilePath string) ([]types.WriteRequest, error) {
	file, err := os.Open(itemFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//
	arWriteRequests := make([]types.WriteRequest, 0, 1000)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strLine := scanner.Text()
		itemObj := make(map[string]interface{})
		err := json.Unmarshal([]byte(strLine), &itemObj)
		if err != nil {
			return nil, err
		}

		itemForPut, _ := attributevalue.MarshalMap(itemObj)
		writeReq := types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: itemForPut,
			},
		}
		arWriteRequests = append(arWriteRequests, writeReq)

		// /*for debug*/
		//jsonOutput, _ := json.MarshalIndent(itemForPut, "", "\t")
		//fmt.Println(string(jsonOutput["PK"]))
		//break
	}
	if err := scanner.Err(); err != nil {
		return nil, err
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
	const strTableName = "quick-photos"
	const strItemFileName = "items.json"
	const maxItemsInBatch = 25

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path = path + string(os.PathSeparator) + strItemFileName
	fmt.Println(path)

	arWriteRequests, err := makeWriteRequests(path)
	lenOfRequests := len(arWriteRequests)
	fmt.Printf("%d items\n", lenOfRequests)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(cfg)

	for startIndex := 0; startIndex < lenOfRequests; startIndex += maxItemsInBatch {
		endIndex := startIndex + maxItemsInBatch
		if endIndex > lenOfRequests {
			endIndex = lenOfRequests
		}

		output, err := BatchWrite(client, strTableName, arWriteRequests[startIndex:endIndex])
		// /*for debug*/
		if output != nil {
			jsonOutput, _ := json.MarshalIndent(output, "", "\t")
			fmt.Println(string(jsonOutput))
		}

		if err != nil {
			fmt.Println("Error!")
			log.Fatal(err)
		}

		time.Sleep(1500 * time.Millisecond)
	}

}
