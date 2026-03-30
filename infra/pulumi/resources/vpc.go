package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type VPCOutput struct {
	VPC            *ec2.Vpc
	PublicSubnets  pulumi.StringArrayOutput
	PrivateSubnets pulumi.StringArrayOutput
}

func CreateVPC(ctx *pulumi.Context) (*VPCOutput, error) {
	// Create VPC with explicit 3 AZs for multi-AZ distribution
	vpc, err := ec2.NewVpc(ctx, "vpc", &ec2.VpcArgs{
		CidrBlock: pulumi.String("10.0.0.0/16"),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("vpc"),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create Internet Gateway
	igw, err := ec2.NewInternetGateway(ctx, "vpc", &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("vpc"),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create public route table
	publicRt, err := ec2.NewRouteTable(ctx, "public-rt", &ec2.RouteTableArgs{
		VpcId: vpc.ID(),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("public-rt"),
		},
	})
	if err != nil {
		return nil, err
	}

	// Add default route to IGW for public route table
	_, err = ec2.NewRoute(ctx, "public-route", &ec2.RouteArgs{
		RouteTableId:         publicRt.ID(),
		DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
		GatewayId:            igw.ID(),
	})
	if err != nil {
		return nil, err
	}

	// Create private route table
	privateRt, err := ec2.NewRouteTable(ctx, "private-rt", &ec2.RouteTableArgs{
		VpcId: vpc.ID(),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("private-rt"),
		},
	})
	if err != nil {
		return nil, err
	}

	// AZs and CIDR blocks
	azs := []string{"sa-east-1a", "sa-east-1b", "sa-east-1c"}
	publicSubnets := make([]pulumi.StringOutput, len(azs))
	privateSubnets := make([]pulumi.StringOutput, len(azs))

	for i, az := range azs {
		// Create public subnet
		pubSubnet, err := ec2.NewSubnet(ctx, fmt.Sprintf("public-subnet-%d", i+1), &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String(fmt.Sprintf("10.0.%d.0/24", i)),
			AvailabilityZone: pulumi.String(az),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(fmt.Sprintf("public-subnet-%d", i+1)),
			},
		})
		if err != nil {
			return nil, err
		}
		publicSubnets[i] = pubSubnet.ID().ToStringOutput()

		// Associate public subnet with public route table
		_, err = ec2.NewRouteTableAssociation(ctx, fmt.Sprintf("public-rta-%d", i+1), &ec2.RouteTableAssociationArgs{
			SubnetId:     pubSubnet.ID(),
			RouteTableId: publicRt.ID(),
		})
		if err != nil {
			return nil, err
		}

		// Create private subnet
		privSubnet, err := ec2.NewSubnet(ctx, fmt.Sprintf("private-subnet-%d", i+1), &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String(fmt.Sprintf("10.0.%d.0/24", 100+i)),
			AvailabilityZone: pulumi.String(az),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(fmt.Sprintf("private-subnet-%d", i+1)),
			},
		})
		if err != nil {
			return nil, err
		}
		privateSubnets[i] = privSubnet.ID().ToStringOutput()

		// Associate private subnet with private route table
		_, err = ec2.NewRouteTableAssociation(ctx, fmt.Sprintf("private-rta-%d", i+1), &ec2.RouteTableAssociationArgs{
			SubnetId:     privSubnet.ID(),
			RouteTableId: privateRt.ID(),
		})
		if err != nil {
			return nil, err
		}

		// Create EIP for NAT Gateway
		eip, err := ec2.NewEip(ctx, fmt.Sprintf("nat-eip-%d", i+1), &ec2.EipArgs{
			Domain: pulumi.String("vpc"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(fmt.Sprintf("nat-eip-%d", i+1)),
			},
		}, pulumi.DependsOn([]pulumi.Resource{igw}))
		if err != nil {
			return nil, err
		}

		// Create NAT Gateway in public subnet
		natGw, err := ec2.NewNatGateway(ctx, fmt.Sprintf("nat-gw-%d", i+1), &ec2.NatGatewayArgs{
			AllocationId: eip.ID(),
			SubnetId:     pubSubnet.ID(),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(fmt.Sprintf("nat-gw-%d", i+1)),
			},
		})
		if err != nil {
			return nil, err
		}

		// Use first NAT Gateway for private route
		if i == 0 {
			// Add route to NAT Gateway for private subnets (only once)
			_, err = ec2.NewRoute(ctx, "private-route", &ec2.RouteArgs{
				RouteTableId:         privateRt.ID(),
				DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
				NatGatewayId:         natGw.ID(),
			}, pulumi.DependsOn([]pulumi.Resource{natGw}))
			if err != nil {
				return nil, err
			}
		}
	}

	return &VPCOutput{
		VPC:            vpc,
		PublicSubnets:  pulumi.ToStringArrayOutput(publicSubnets),
		PrivateSubnets: pulumi.ToStringArrayOutput(privateSubnets),
	}, nil
}
