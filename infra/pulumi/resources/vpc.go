package resources

import (
	awsxec2 "github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type VPCOutput struct {
	VPC            *awsxec2.Vpc
	PublicSubnets  pulumi.StringArrayOutput
	PrivateSubnets pulumi.StringArrayOutput
}

func CreateVPC(ctx *pulumi.Context) (*VPCOutput, error) {
	// Create VPC with explicit 3 AZs to force multi-AZ distribution
	// This ensures RDS and ALB have subnets across multiple availability zones
	vpc, err := awsxec2.NewVpc(ctx, "vpc", &awsxec2.VpcArgs{
		CidrBlock: pulumi.StringRef("10.0.0.0/16"),
		AvailabilityZoneNames: []string{
			"us-east-1a",
			"us-east-1b",
			"us-east-1c",
		},
		NatGateways: &awsxec2.NatGatewayConfigurationArgs{
			Strategy: awsxec2.NatGatewayStrategyOnePerAz,
		},
	})

	if err != nil {
		return nil, err
	}

	return &VPCOutput{
		VPC:            vpc,
		PublicSubnets:  vpc.PublicSubnetIds,
		PrivateSubnets: vpc.PrivateSubnetIds,
	}, nil
}
