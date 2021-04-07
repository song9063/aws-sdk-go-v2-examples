package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {
	const instanceId = ""
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := ec2.NewFromConfig(cfg)

	// no filter
	//output, err := client.DescribeInstances(context.TODO(), nil)

	// filter
	output, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []string{
					"running",
					"pending",
					"stopped",
				},
			},
		},
		InstanceIds: []string{
			instanceId,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	jsonOutput, _ := json.MarshalIndent(output, "", "\t")
	fmt.Println(string(jsonOutput))
}
