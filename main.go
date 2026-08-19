package main

import (
	"fmt"
	"os"

	awsconfig "ssm-cli/awsConfig"
	ec2 "ssm-cli/cmd/ec2"
	fzf "ssm-cli/cmd/fzf"
	ssmconnect "ssm-cli/cmd/ssmConnect"
)


func main() {
	region := ""
	if len(os.Args) > 1 {
		region = os.Args[1]
	}

	cfg, err := awsconfig.SetupAWSConfig(region)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to set up AWS config:", err)
		os.Exit(1)
	}

	ec2InstancesId, ec2InstanceName, err := ec2.ListEC2Instances(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to list EC2 instances:", err)
		os.Exit(1)
	}

	fzfSelected, err := fzf.FzfSelect(ec2InstancesId, ec2InstanceName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to select EC2 instance:", err)
		os.Exit(1)
	}

	fmt.Println("Selected EC2 Instance:", fzfSelected)

	err = ssmconnect.ConnectToEC2Instance(cfg, fzfSelected)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to EC2 instance:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to EC2 Instance:", fzfSelected)


}

