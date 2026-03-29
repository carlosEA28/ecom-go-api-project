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
	// Security Group para RDS
	rdsSecurityGroup, err := awsec2.NewSecurityGroup(ctx, "rds-security-group", &awsec2.SecurityGroupArgs{
		VpcId: vpcID,
		Ingress: awsec2.SecurityGroupIngressArray{
			&awsec2.SecurityGroupIngressArgs{
				FromPort: pulumi.Int(5432),
				ToPort:   pulumi.Int(5432),
				Protocol: pulumi.String("tcp"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("10.0.0.0/16"), // Allow from VPC
				},
				Description: pulumi.StringPtr("Allow PostgreSQL from ECS"),
			},
		},
		Egress: awsec2.SecurityGroupEgressArray{
			// RDS doesn't need outbound access
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String("ecom-rds-sg"),
		},
	})
	if err != nil {
		return nil, err
	}

	// Security Group para ECS
	ecsSecurityGroup, err := awsec2.NewSecurityGroup(ctx, "ecs-security-group", &awsec2.SecurityGroupArgs{
		VpcId: vpcID,
		Ingress: awsec2.SecurityGroupIngressArray{
			&awsec2.SecurityGroupIngressArgs{
				FromPort: pulumi.Int(8080),
				ToPort:   pulumi.Int(8080),
				Protocol: pulumi.String("tcp"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("10.0.0.0/16"), // Allow from LB
				},
				Description: pulumi.StringPtr("Allow application traffic"),
			},
		},
		Egress: awsec2.SecurityGroupEgressArray{
			// ECS to HTTPS (external APIs, package downloads)
			&awsec2.SecurityGroupEgressArgs{
				FromPort: pulumi.Int(443),
				ToPort:   pulumi.Int(443),
				Protocol: pulumi.String("tcp"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
				Description: pulumi.StringPtr("Allow HTTPS outbound"),
			},
			// ECS to DNS
			&awsec2.SecurityGroupEgressArgs{
				FromPort: pulumi.Int(53),
				ToPort:   pulumi.Int(53),
				Protocol: pulumi.String("udp"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
				Description: pulumi.StringPtr("Allow DNS queries"),
			},
			// ECS to RDS
			&awsec2.SecurityGroupEgressArgs{
				FromPort: pulumi.Int(5432),
				ToPort:   pulumi.Int(5432),
				Protocol: pulumi.String("tcp"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("10.0.0.0/16"),
				},
				Description: pulumi.StringPtr("Allow PostgreSQL"),
			},
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String("ecom-ecs-sg"),
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
				Description: pulumi.StringPtr("Allow HTTP from internet"),
			},
			&awsec2.SecurityGroupIngressArgs{
				FromPort: pulumi.Int(443),
				ToPort:   pulumi.Int(443),
				Protocol: pulumi.String("tcp"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
				Description: pulumi.StringPtr("Allow HTTPS from internet"),
			},
		},
		Egress: awsec2.SecurityGroupEgressArray{
			// LB to ECS
			&awsec2.SecurityGroupEgressArgs{
				FromPort: pulumi.Int(8080),
				ToPort:   pulumi.Int(8080),
				Protocol: pulumi.String("tcp"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("10.0.0.0/16"),
				},
				Description: pulumi.StringPtr("Allow to ECS application"),
			},
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String("ecom-lb-sg"),
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
