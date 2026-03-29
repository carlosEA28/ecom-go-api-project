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
	vpc, err := awsxec2.NewVpc(ctx, "vpc", &awsxec2.VpcArgs{
		CidrBlock: pulumi.StringRef("10.0.0.0/24"),
		SubnetSpecs: []awsxec2.SubnetSpecArgs{
			{
				Name:     pulumi.StringRef("public-subnet-loadBalancer"),
				Type:     awsxec2.SubnetTypePublic,
				CidrMask: pulumi.IntRef(22),
			},
			{
				Name:     pulumi.StringRef("ECS-cluster-subnetprivate-subnet"),
				Type:     awsxec2.SubnetTypePrivate,
				CidrMask: pulumi.IntRef(22),
			},
			{
				Name:     pulumi.StringRef("RDS-subnetprivate-subnet"),
				Type:     awsxec2.SubnetTypePrivate,
				CidrMask: pulumi.IntRef(22),
			},
		},
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
