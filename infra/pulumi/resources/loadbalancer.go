package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type LoadBalancerOutput struct {
	LoadBalancer *lb.ApplicationLoadBalancer
	DNSName      pulumi.StringOutput
	TargetGroup  pulumi.Output
}

func CreateLoadBalancer(ctx *pulumi.Context, publicSubnetID pulumi.Output, lbSecurityGroupID pulumi.StringOutput) (*LoadBalancerOutput, error) {
	cfg := config.New(ctx, "")
	environment := cfg.Get("environment")

	// Create the load balancer with HTTP listener
	// In production, integrate with ACM certificate for HTTPS
	loadBalancer, err := lb.NewApplicationLoadBalancer(ctx, "lb", &lb.ApplicationLoadBalancerArgs{
		Listener: &lb.ListenerArgs{
			Port:     pulumi.Int(80),
			Protocol: pulumi.String("HTTP"),
			// TODO: Add HTTPS listener when ACM certificate is available
			// Configure with: CertificateArn from Pulumi config
		},
		SubnetIds: pulumi.StringArray{
			publicSubnetID.ApplyT(func(id interface{}) string {
				return id.(string)
			}).(pulumi.StringOutput),
		},
		SecurityGroups: pulumi.StringArray{
			lbSecurityGroupID,
		},
		Tags: pulumi.StringMap{
			"Name":        pulumi.String("ecom-lb"),
			"Environment": pulumi.String(environment),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create load balancer: %w", err)
	}

	// Log for documentation
	ctx.Log.Info(fmt.Sprintf("Load Balancer created for %s environment", environment), nil)
	if environment == "prod" || environment == "staging" {
		ctx.Log.Warn("Consider adding HTTPS listener with ACM certificate for production", nil)
	}

	return &LoadBalancerOutput{
		LoadBalancer: loadBalancer,
		DNSName:      loadBalancer.LoadBalancer.DnsName(),
		TargetGroup:  loadBalancer.DefaultTargetGroup,
	}, nil
}
