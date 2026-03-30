package resources

import (
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type LoadBalancerOutput struct {
	LoadBalancer   *lb.LoadBalancer
	TargetGroup    *lb.TargetGroup
	DNSName        pulumi.StringOutput
	TargetGroupArn pulumi.StringOutput
}

func CreateLoadBalancer(ctx *pulumi.Context, vpcID pulumi.StringOutput, publicSubnetIDs pulumi.StringArrayOutput, lbSecurityGroupID pulumi.StringOutput) (*LoadBalancerOutput, error) {
	// Create Load Balancer
	loadBalancer, err := lb.NewLoadBalancer(ctx, "lb", &lb.LoadBalancerArgs{
		Internal:         pulumi.Bool(false),
		LoadBalancerType: pulumi.String("application"),
		SecurityGroups:   pulumi.StringArray{lbSecurityGroupID},
		Subnets:          publicSubnetIDs,
		Tags: pulumi.StringMap{
			"Name": pulumi.String("ecom-api-lb"),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create Target Group
	targetGroup, err2 := lb.NewTargetGroup(ctx, "lb", &lb.TargetGroupArgs{
		Port:       pulumi.Int(8080),
		Protocol:   pulumi.String("HTTP"),
		VpcId:      vpcID,
		TargetType: pulumi.String("ip"),
		HealthCheck: &lb.TargetGroupHealthCheckArgs{
			Path:     pulumi.String("/health"),
			Protocol: pulumi.String("HTTP"),
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String("ecom-api-tg"),
		},
	})
	if err2 != nil {
		return nil, err2
	}

	// Create Listener
	_, err3 := lb.NewListener(ctx, "lb", &lb.ListenerArgs{
		LoadBalancerArn: loadBalancer.Arn,
		Port:            pulumi.Int(80),
		Protocol:        pulumi.String("HTTP"),
		DefaultActions: lb.ListenerDefaultActionArray{
			&lb.ListenerDefaultActionArgs{
				Type:           pulumi.String("forward"),
				TargetGroupArn: targetGroup.Arn,
			},
		},
	})
	if err3 != nil {
		return nil, err3
	}

	return &LoadBalancerOutput{
		LoadBalancer:   loadBalancer,
		TargetGroup:    targetGroup,
		DNSName:        loadBalancer.DnsName,
		TargetGroupArn: targetGroup.Arn,
	}, nil
}
