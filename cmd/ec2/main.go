package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func ListEC2Instances(config aws.Config) ([]string, []string, error) {

	ec2Client := ec2.NewFromConfig(config)

	listInstancesInput := &ec2.DescribeInstancesInput{}

	listInstancesOutput, err := ec2Client.DescribeInstances(context.Background(), listInstancesInput)
	
	if err != nil {
		return nil, nil, err
	}

	var instanceIDs []string
	var instanceNames []string

	for _, reservation := range listInstancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" && tag.Value != nil {
					instanceNames = append(instanceNames, *tag.Value)
				}		

			}
		}
	}

	return instanceIDs, instanceNames, nil
}