package resources

import (
	awsec2 "github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SecurityGroupOutput struct {
	ECSSecurityGroup *awsec2.SecurityGroup
	RDSSecurityGroup *awsec2.SecurityGroup
	LBSecurityGroup  *awsec2.SecurityGroup
}

func CreateSecurityGroups(ctx *pulumi.Context, vpcID pulumi.StringOutput) (*SecurityGroupOutput, error) {
	// Security Group para ECS
	ecsSecurityGroup, err := awsec2.NewSecurityGroup(ctx, "securityGroup", &awsec2.SecurityGroupArgs{
		VpcId: vpcID,
		Ingress: awsec2.SecurityGroupIngressArray{
			&awsec2.SecurityGroupIngressArgs{
				FromPort: pulumi.Int(8080),
				ToPort:   pulumi.Int(8080),
				Protocol: pulumi.String("tcp"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("10.0.0.0/24"),
				},
			},
		},
		Egress: awsec2.SecurityGroupEgressArray{
			&awsec2.SecurityGroupEgressArgs{
				FromPort: pulumi.Int(0),
				ToPort:   pulumi.Int(0),
				Protocol: pulumi.String("-1"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
				Ipv6CidrBlocks: pulumi.StringArray{
					pulumi.String("::/0"),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Security Group para RDS
	rdsSecurityGroup, err := awsec2.NewSecurityGroup(ctx, "rds-security-group", &awsec2.SecurityGroupArgs{
		VpcId: vpcID,
		Ingress: awsec2.SecurityGroupIngressArray{
			&awsec2.SecurityGroupIngressArgs{
				FromPort: pulumi.Int(5432),
				ToPort:   pulumi.Int(5432),
				Protocol: pulumi.String("tcp"),
				SecurityGroups: pulumi.StringArray{
					ecsSecurityGroup.ID(),
				},
			},
		},
		Egress: awsec2.SecurityGroupEgressArray{
			&awsec2.SecurityGroupEgressArgs{
				FromPort: pulumi.Int(0),
				ToPort:   pulumi.Int(0),
				Protocol: pulumi.String("-1"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Security Group para Load Balancer
	lbSecurityGroup, err := awsec2.NewSecurityGroup(ctx, "lb-security-group", &awsec2.SecurityGroupArgs{
		VpcId: vpcID,
		Ingress: awsec2.SecurityGroupIngressArray{
			&awsec2.SecurityGroupIngressArgs{
				FromPort: pulumi.Int(80),
				ToPort:   pulumi.Int(80),
				Protocol: pulumi.String("tcp"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
			},
		},
		Egress: awsec2.SecurityGroupEgressArray{
			&awsec2.SecurityGroupEgressArgs{
				FromPort: pulumi.Int(0),
				ToPort:   pulumi.Int(0),
				Protocol: pulumi.String("-1"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &SecurityGroupOutput{
		ECSSecurityGroup: ecsSecurityGroup,
		RDSSecurityGroup: rdsSecurityGroup,
		LBSecurityGroup:  lbSecurityGroup,
	}, nil
}
