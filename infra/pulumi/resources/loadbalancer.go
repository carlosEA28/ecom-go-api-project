package resources

import (
	"github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type LoadBalancerOutput struct {
	LoadBalancer *lb.ApplicationLoadBalancer
	DNSName      pulumi.StringOutput
}

func CreateLoadBalancer(ctx *pulumi.Context, publicSubnetIDs pulumi.StringArrayOutput, lbSecurityGroupID pulumi.StringOutput) (*LoadBalancerOutput, error) {
	loadBalancer, err := lb.NewApplicationLoadBalancer(ctx, "lb", &lb.ApplicationLoadBalancerArgs{
		Listener: &lb.ListenerArgs{
			Port:     pulumi.Int(80),
			Protocol: pulumi.String("HTTP"),
		},
		SubnetIds: publicSubnetIDs,
		SecurityGroups: pulumi.StringArray{
			lbSecurityGroupID,
		},
	})

	if err != nil {
		return nil, err
	}

	return &LoadBalancerOutput{
		LoadBalancer: loadBalancer,
		DNSName:      loadBalancer.LoadBalancer.DnsName(),
	}, nil
}
