package main

import (
	"fmt"

	"project/resources"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		vpcOutput, err := resources.CreateVPC(ctx)
		if err != nil {
			return fmt.Errorf("erro ao criar VPC: %w", err)
		}

		securityGroups, err := resources.CreateSecurityGroups(ctx, vpcOutput.VPC.ID().ToStringOutput())
		if err != nil {
			return fmt.Errorf("erro ao criar Security Groups: %w", err)
		}

		ecrOutput, err := resources.CreateECR(ctx)
		if err != nil {
			return fmt.Errorf("erro ao criar ECR: %w", err)
		}

		// Use all public subnets for load balancer (multi-AZ)
		// Use all private subnets for RDS and ECS (multi-AZ)

		loadBalancerOutput, err := resources.CreateLoadBalancer(
			ctx,
			vpcOutput.VPC.ID().ToStringOutput(),
			vpcOutput.PublicSubnets,
			securityGroups.LBSecurityGroup.ID().ToStringOutput(),
		)
		if err != nil {
			return fmt.Errorf("erro ao criar Load Balancer: %w", err)
		}

		rdsOutput, err := resources.CreateRDS(
			ctx,
			vpcOutput.PrivateSubnets,
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
			vpcOutput.PrivateSubnets,
			securityGroups.ECSSecurityGroup.ID().ToStringOutput(),
			loadBalancerOutput,
		)
		if err != nil {
			return fmt.Errorf("erro ao criar ECS Fargate Service: %w", err)
		}

		// Exports
		ctx.Export("vpcId", vpcOutput.VPC.ID())

		ctx.Export("ecrRepositoryUrl", resources.GetECRRepositoryUrl(ecrOutput))
		ctx.Export("ecrImageUri", ecrOutput.ImageURI)

		ctx.Export("loadBalancerDns", loadBalancerOutput.DNSName)

		ctx.Export("rdsEndpoint", rdsOutput.Endpoint)
		ctx.Export("rdsPort", rdsOutput.Port)
		ctx.Export("rdsDatabase", rdsOutput.DatabaseName)
		ctx.Export("rdsUsername", rdsOutput.Username)

		ctx.Export("ecsClusterArn", ecsClusterOutput.Cluster.Arn)
		ctx.Export("ecsServiceArn", ecsServiceOutput.Service.Arn)

		return nil

	})
}
