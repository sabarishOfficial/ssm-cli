package ssmConnect

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func ConnectToEC2Instance(config aws.Config, instanceID string) error {
	if instanceID == "" {
		return fmt.Errorf("instance ID cannot be empty")
	}

	ssmClient := ssm.NewFromConfig(config)

	output, err := ssmClient.StartSession(context.Background(), &ssm.StartSessionInput{
		Target: aws.String(instanceID),
	})
	if err != nil {
		return fmt.Errorf("failed to start SSM session: %w", err)
	}

	sessionJSON, err := json.Marshal(output)
	if err != nil {
		return fmt.Errorf("failed to marshal session output: %w", err)
	}

	paramsJSON, err := json.Marshal(map[string]string{"Target": instanceID})
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	endpoint := fmt.Sprintf("https://ssm.%s.amazonaws.com", config.Region)

	// session-manager-plugin must be installed on the host
	cmd := exec.Command(
		"session-manager-plugin",
		string(sessionJSON),
		config.Region,
		"StartSession",
		"",
		string(paramsJSON),
		endpoint,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}