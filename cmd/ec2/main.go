package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	ami "github.com/Jkarage/AMI-update/internal"
	"github.com/ardanlabs/conf/v3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func run() error {
	cfg := struct {
		Region       string `conf:"default:eu-west-3"`
		InstanceID   string `conf:"default:i-04493d3e5d3001e32"`
		Script       string `conf:"default:./internal/script.txt"`
		AWSACCESSID  string `conf:""`
		AWSSECRETKEY string `conf:""`
	}{}

	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return err
	}

	awsConf, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.Region))
	if err != nil {
		return err
	}

	cmds, err := ami.ReadScript(cfg.Script)
	if err != nil {
		return err
	}

	ssmClient := ssm.NewFromConfig(awsConf)
	targets := []types.Target{
		{
			Key:    aws.String("instanceids"),
			Values: []string{cfg.InstanceID},
		},
	}

	// Create a document specifying the command to be executed
	docName := "AWS-RunShellScript"
	docInput := &ssm.SendCommandInput{
		DocumentName: &docName,
		Targets:      targets,
		Parameters: map[string][]string{
			"commands": cmds,
		},
	}

	out, err := ssmClient.SendCommand(context.TODO(), docInput)
	if err != nil {
		return err
	}

	outIn := ssm.GetCommandInvocationInput{
		CommandId:  out.Command.CommandId,
		InstanceId: &cfg.InstanceID,
	}

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		invocationOutput, err := ssmClient.GetCommandInvocation(context.TODO(), &outIn)
		if err != nil {
			fmt.Println("ERRor getting output", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if invocationOutput.Status == types.CommandInvocationStatusSuccess {
			fmt.Println(*invocationOutput.StandardOutputContent)
			break
		} else {
			fmt.Println("Command execution not yet completed. Waiting...")
			time.Sleep(2 * time.Second)
			continue
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
