package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	stypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func main() {
	cfg, err := Config("eu-west-3")
	if err != nil {
		log.Fatal(err.Error())
	}

	ec2Client := ec2.NewFromConfig(cfg)

	ami, err := GetAMI("gpt-template", cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ami)

	instanceId, err := LaunchInstance(ami, ec2Client)
	if err != nil {
		log.Fatal(err)
	}

	err = gitPull(cfg, instanceId, "sarufi-vector-index")
	if err != nil {
		log.Fatal(err)
	}
}

// Config populates aws Config with aws required properties
// expecting the  credentials to be in the environment.
func Config(region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

// gitPull runs git pull command in an ec2 instance.
// This is for updating the running service files
func gitPull(cfg aws.Config, instance, dir string) error {
	client := ssm.NewFromConfig(cfg)

	// Specify the command to be executed
	command := []string{
		"touch /home/ubuntu/a.txt",
		"echo 'Hello from SSM command!' > /home/ubuntu/a.txt",
		"ls -lah",
	}

	// Specify the target instance
	targets := []stypes.Target{
		{
			Key:    aws.String("instanceids"),
			Values: []string{instance},
		},
	}

	// Create a document specifying the command to be executed
	docName := "AWS-RunShellScript"
	docInput := &ssm.SendCommandInput{
		DocumentName: &docName,
		Targets:      targets,
		Parameters: map[string][]string{
			"commands": command,
		},
	}

	// Send the command
	output, err := client.SendCommand(context.TODO(), docInput)
	if err != nil {
		fmt.Println("Error sending command:", err)
		os.Exit(1)
	}

	// Print command ID
	fmt.Println("Command ID:", *output.Command.CommandId)

	return nil
}

// GetAMI gets instance info from a launching template
func GetAMI(lt string, cfg aws.Config) (string, error) {
	ec2Client := ec2.NewFromConfig(cfg)

	params := ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateName: &lt,
	}

	ltOuptut, err := ec2Client.DescribeLaunchTemplateVersions(context.TODO(), &params)
	if err != nil {
		return "", err
	}

	if len(ltOuptut.LaunchTemplateVersions) == 0 {
		return "", fmt.Errorf("the launch template %s not found", lt)
	}

	return *ltOuptut.LaunchTemplateVersions[0].LaunchTemplateData.ImageId, nil
}

// Launch an instance using the latest ami
// retrieve the new instance id
func LaunchInstance(ami string, client *ec2.Client) (string, error) {
	var count int32 = 1

	input := ec2.RunInstancesInput{
		ImageId:      &ami,
		MinCount:     &count,
		MaxCount:     &count,
		InstanceType: types.InstanceTypeC5Large,
	}

	resp, err := client.RunInstances(context.TODO(), &input)
	if err != nil {
		return "", err
	}

	if len(resp.Instances) == 0 {
		return "", fmt.Errorf("no instance launched")
	}

	return *resp.Instances[0].InstanceId, nil
}
