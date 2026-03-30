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
	// Create VPC with default subnet strategy which spreads across multiple AZs
	vpc, err := awsxec2.NewVpc(ctx, "vpc", &awsxec2.VpcArgs{
		CidrBlock: pulumi.StringRef("10.0.0.0/16"),
		// Omit SubnetSpecs to use default strategy (which creates subnets in multiple AZs)
		NatGateways: &awsxec2.NatGatewayConfigurationArgs{
			Strategy: awsxec2.NatGatewayStrategySingle,
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
