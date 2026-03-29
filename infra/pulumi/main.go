package main

import (
	"fmt"

	"github.com/carlosEA28/ecom/infra/pulumi/resources"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		vpcOutput, err := resources.CreateVPC(ctx)
		if err != nil {
			return fmt.Errorf("erro ao criar VPC: %w", err)
		}

		securityGroups, err := resources.CreateSecurityGroups(ctx, vpcOutput.VPC.VpcId)
		if err != nil {
			return fmt.Errorf("erro ao criar Security Groups: %w", err)
		}

		ecrOutput, err := resources.CreateECR(ctx)
		if err != nil {
			return fmt.Errorf("erro ao criar ECR: %w", err)
		}

		publicSubnetID := vpcOutput.PublicSubnets.Index(pulumi.Int(0))
		ecsSubnetID := vpcOutput.PrivateSubnets.Index(pulumi.Int(0))
		rdsSubnetID := vpcOutput.PrivateSubnets.Index(pulumi.Int(1))

		loadBalancerOutput, err := resources.CreateLoadBalancer(
			ctx,
			publicSubnetID,
			securityGroups.LBSecurityGroup.ID(),
		)
		if err != nil {
			return fmt.Errorf("erro ao criar Load Balancer: %w", err)
		}

		rdsOutput, err := resources.CreateRDS(
			ctx,
			rdsSubnetID,
			securityGroups.RDSSecurityGroup.ID(),
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
			securityGroups.ECSSecurityGroup.ID(),
			loadBalancerOutput.LoadBalancer,
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
		ctx.Export("ecsServiceName", ecsServiceOutput.Service.ServiceArn)

		return nil
	})
}
