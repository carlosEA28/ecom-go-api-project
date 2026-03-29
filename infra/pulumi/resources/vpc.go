package resources

import (
	"fmt"

	awsxec2 "github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type VPCOutput struct {
	VPC            *awsxec2.Vpc
	PublicSubnets  pulumi.StringArrayOutput
	PrivateSubnets pulumi.StringArrayOutput
}

func CreateVPC(ctx *pulumi.Context) (*VPCOutput, error) {
	cfg := config.New(ctx, "")
	environment := cfg.Get("environment")

	// Use OnePerAz strategy for production/staging (high availability)
	// Use Single for dev (cost optimization)
	natGatewayStrategy := awsxec2.NatGatewayStrategySingle
	if environment == "prod" || environment == "staging" {
		natGatewayStrategy = awsxec2.NatGatewayStrategyOnePerAz
	}

	vpc, err := awsxec2.NewVpc(ctx, "vpc", &awsxec2.VpcArgs{
		CidrBlock: pulumi.StringRef("10.0.0.0/16"),
		SubnetSpecs: []awsxec2.SubnetSpecArgs{
			{
				Name:     pulumi.StringRef("public-subnet-loadBalancer"),
				Type:     awsxec2.SubnetTypePublic,
				CidrMask: pulumi.IntRef(24),
			},
			{
				Name:     pulumi.StringRef("ecs-private-subnet"),
				Type:     awsxec2.SubnetTypePrivate,
				CidrMask: pulumi.IntRef(24),
			},
			{
				Name:     pulumi.StringRef("rds-private-subnet"),
				Type:     awsxec2.SubnetTypePrivate,
				CidrMask: pulumi.IntRef(24),
			},
		},
		NatGateways: &awsxec2.NatGatewayConfigurationArgs{
			Strategy: natGatewayStrategy,
		},
		Tags: pulumi.StringMap{
			"Name":        pulumi.String("ecom-vpc"),
			"Environment": pulumi.String(environment),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create VPC: %w", err)
	}

	strategyName := "Single"
	if natGatewayStrategy == awsxec2.NatGatewayStrategyOnePerAz {
		strategyName = "OnePerAz"
	}
	ctx.Log.Info(fmt.Sprintf("VPC created with NAT strategy: %s for %s environment", strategyName, environment), nil)

	// TODO: Enable VPC Flow Logs for network monitoring
	// VPC Flow Logs help with:
	// - Network troubleshooting
	// - Security analysis
	// - Compliance requirements
	// Configuration example for CloudWatch:
	// - Log destination: CloudWatch Logs group
	// - Traffic type: ACCEPT and REJECT
	// - Retention: 30 days for staging/prod, 7 days for dev

	return &VPCOutput{
		VPC:            vpc,
		PublicSubnets:  vpc.PublicSubnetIds,
		PrivateSubnets: vpc.PrivateSubnetIds,
	}, nil
}
