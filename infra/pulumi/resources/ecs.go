package resources

import (
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecs"
	awsxecs "github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/ecs"
	ecsx "github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/ecs"
	"github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ECSClusterOutput struct {
	Cluster *ecs.Cluster
}

type ECSServiceOutput struct {
	Service *awsxecs.FargateService
}

func CreateECSCluster(ctx *pulumi.Context) (*ECSClusterOutput, error) {
	cluster, err := ecs.NewCluster(ctx, "cluster", nil)
	if err != nil {
		return nil, err
	}

	return &ECSClusterOutput{
		Cluster: cluster,
	}, nil
}

func CreateECSFargateService(
	ctx *pulumi.Context,
	cluster *ecs.Cluster,
	imageURI pulumi.StringOutput,
	ecsSubnetIDs pulumi.StringArrayOutput,
	ecsSecurityGroupID pulumi.StringOutput,
	loadBalancer *lb.ApplicationLoadBalancer,
) (*ECSServiceOutput, error) {
	service, err := awsxecs.NewFargateService(ctx, "service", &awsxecs.FargateServiceArgs{
		Cluster:                       cluster.Arn,
		HealthCheckGracePeriodSeconds: pulumi.Int(300),
		NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
			Subnets: ecsSubnetIDs,
			SecurityGroups: pulumi.StringArray{
				ecsSecurityGroupID,
			},
		},
		DesiredCount: pulumi.Int(1),
		TaskDefinitionArgs: &awsxecs.FargateServiceTaskDefinitionArgs{
			Container: &awsxecs.TaskDefinitionContainerDefinitionArgs{
				Image:  imageURI,
				Cpu:    pulumi.Int(256),
				Memory: pulumi.Int(512),
				PortMappings: ecsx.TaskDefinitionPortMappingArray{
					ecsx.TaskDefinitionPortMappingArgs{
						ContainerPort: pulumi.Int(8080),
						HostPort:      pulumi.Int(8080),
						TargetGroup:   loadBalancer.DefaultTargetGroup,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &ECSServiceOutput{
		Service: service,
	}, nil
}
