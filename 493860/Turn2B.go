package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

// CloudProvider interface defines the method to get the cloud provider name
type CloudProvider interface {
	Name() string
	CollectMetrics(ctx context.Context, results chan<- result) error
}

type gcpProvider struct{}

func (p *gcpProvider) Name() string {
	return "GCP"
}

func (p *gcpProvider) CollectMetrics(ctx context.Context, results chan<- result) error {
	// GCP implementation to collect metrics
	return nil
}

type awsProvider struct{}

func (p *awsProvider) Name() string {
	return "AWS"
}

func (p *awsProvider) CollectMetrics(ctx context.Context, results chan<- result) error {
	// AWS implementation to collect metrics
	// Set AWS region from metadata
	region, err := ec2metadata.New(ctx).Region()
	if err != nil {
		return err
	}

	// Create an AWS session
	cfg, err := aws.LoadDefaultConfig(ctx, aws.WithRegion(region))
	if err != nil {
		return err
	}

	// Create an EC2 service client
	svc := ec2.NewFromConfig(cfg)

	// Describe the instance to get its ID
	describeInstancesOutput, err := svc.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return err
	}

	if len(describeInstancesOutput.Reservations) == 0 || len(describeInstancesOutput.Reservations[0].Instances) == 0 {
		return fmt.Errorf("no instance found")
	}

	instanceID := *describeInstancesOutput.Reservations[0].Instances[0].InstanceId

	// Describe the instance metrics
	describeInstanceMetricsOutput, err := svc.DescribeInstanceMetrics(ctx, &ec2.DescribeInstanceMetricsInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return err
	}

	for _, metricData := range describeInstanceMetricsOutput.MetricDataResults {
		for _, value := range metricData.Values {
			latency, err := time.ParseDuration(fmt.Sprintf("%.2fms", *value))
			if err != nil {
				continue
			}
			results <- result{latency: latency}
		}
	}

	return nil
}

type azureProvider struct{}

func (p *azureProvider) Name() string {
	return "Azure"
}

func (p *azureProvider) CollectMetrics(ctx context.Context, results chan<- result) error {
	// Azure implementation to collect metrics
	// Set Azure region from metadata
	region, err := metadata.Get("instance/location")
	if err != nil {
		return err
	}

	// Create an Azure subscription credential
	cred, err := azcore.NewDefaultAzureCredential(nil)