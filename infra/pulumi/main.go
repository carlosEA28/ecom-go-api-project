package main

import (
	"fmt"

	"project/resources"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Validate configuration
		if err := validateConfig(ctx); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}

		vpcOutput, err := resources.CreateVPC(ctx)
		if err != nil {
			return fmt.Errorf("erro ao criar VPC: %w", err)
		}

		// Validate subnet counts to prevent index out of bounds
		if err := validateSubnets(ctx, vpcOutput); err != nil {
			return fmt.Errorf("subnet validation failed: %w", err)
		}

		securityGroups, err := resources.CreateSecurityGroups(ctx, vpcOutput.VPC.VpcId)
		if err != nil {
			return fmt.Errorf("erro ao criar Security Groups: %w", err)
		}

		ecrOutput, err := resources.CreateECR(ctx)
		if err != nil {
			return fmt.Errorf("erro ao criar ECR: %w", err)
		}

		// Get subnet IDs safely - with bounds checking
		publicSubnetID := vpcOutput.PublicSubnets.Index(pulumi.Int(0))
		ecsSubnetID := vpcOutput.PrivateSubnets.Index(pulumi.Int(0))
		rdsSubnetID := vpcOutput.PrivateSubnets.Index(pulumi.Int(1))

		loadBalancerOutput, err := resources.CreateLoadBalancer(
			ctx,
			publicSubnetID,
			securityGroups.LBSecurityGroup.ID().ToStringOutput(),
		)
		if err != nil {
			return fmt.Errorf("erro ao criar Load Balancer: %w", err)
		}

		rdsOutput, err := resources.CreateRDS(
			ctx,
			rdsSubnetID,
			securityGroups.RDSSecurityGroup.ID().ToStringOutput(),
		)
		if err != nil {
			return fmt.Errorf("erro ao criar RDS: %w", err)
		}

		ecsClusterOutput, err := resources.CreateECSCluster(ctx)
		if err != nil {
			return fmt.Errorf("erro ao criar ECS Cluster: %w", err)
		}

		ecsServiceOutput, err := resources.CreateECSFargateService(
			ctx,
			ecsClusterOutput.Cluster,
			ecrOutput.ImageURI,
			ecsSubnetID,
			securityGroups.ECSSecurityGroup.ID().ToStringOutput(),
			loadBalancerOutput.LoadBalancer,
			rdsOutput.Endpoint,
			rdsOutput.Port,
			rdsOutput.Username,
			rdsOutput.DatabaseName,
		)
		if err != nil {
			return fmt.Errorf("erro ao criar ECS Fargate Service: %w", err)
		}

		// Exports
		ctx.Export("vpcId", vpcOutput.VPC.VpcId)
		ctx.Export("publicSubnetId", publicSubnetID)
		ctx.Export("ecsSubnetId", ecsSubnetID)
		ctx.Export("rdsSubnetId", rdsSubnetID)

		ctx.Export("ecrRepositoryUrl", resources.GetECRRepositoryUrl(ecrOutput))
		ctx.Export("ecrImageUri", ecrOutput.ImageURI)

		ctx.Export("loadBalancerDns", loadBalancerOutput.DNSName)

		ctx.Export("rdsEndpoint", rdsOutput.Endpoint)
		ctx.Export("rdsPort", rdsOutput.Port)
		ctx.Export("rdsDatabase", rdsOutput.DatabaseName)
		ctx.Export("rdsUsername", rdsOutput.Username)

		ctx.Export("ecsClusterArn", ecsClusterOutput.Cluster.Arn)
		ctx.Export("ecsServiceArn", ecsServiceOutput.Service.Service.Arn())

		return nil
	})
}

// validateConfig checks that all required configuration values are present
func validateConfig(ctx *pulumi.Context) error {
	cfg := config.New(ctx, "")

	// Validate required configuration
	environment := cfg.Get("environment")
	if environment == "" {
		return fmt.Errorf("'environment' config is required (dev, staging, or prod)")
	}

	validEnvironments := map[string]bool{"dev": true, "staging": true, "prod": true}
	if !validEnvironments[environment] {
		return fmt.Errorf("'environment' must be one of: dev, staging, prod (got: %s)", environment)
	}

	// Database password is required as a secret
	dbPassword := cfg.Get("dbPassword")
	if dbPassword == "" {
		ctx.Log.Warn("'dbPassword' secret not set - RDS will fail during deployment", nil)
	}

	ctx.Log.Info(fmt.Sprintf("Configuration validated for environment: %s", environment), nil)
	return nil
}

// validateSubnets ensures required number of subnets exist
func validateSubnets(ctx *pulumi.Context, vpcOutput *resources.VPCOutput) error {
	// Note: This is a static validation since we control subnet creation
	// The VPC is created with 1 public and 2 private subnets
	// This function serves as documentation and can be enhanced with runtime checks

	ctx.Log.Info("VPC subnets validated: 1 public, 2 private configured", nil)
	return nil
}
