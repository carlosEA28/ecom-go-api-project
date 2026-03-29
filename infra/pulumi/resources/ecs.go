package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecs"
	awsxecs "github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/ecs"
	ecsx "github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/ecs"
	"github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type ECSClusterOutput struct {
	Cluster *ecs.Cluster
}

type ECSServiceOutput struct {
	Service *awsxecs.FargateService
}

func CreateECSCluster(ctx *pulumi.Context) (*ECSClusterOutput, error) {
	cfg := config.New(ctx, "")
	environment := cfg.Get("environment")

	cluster, err := ecs.NewCluster(ctx, "cluster", &ecs.ClusterArgs{
		Tags: pulumi.StringMap{
			"Name":        pulumi.String("ecom-cluster"),
			"Environment": pulumi.String(environment),
		},
		// Enable container insights for monitoring
		Settings: ecs.ClusterSettingArray{
			&ecs.ClusterSettingArgs{
				Name:  pulumi.String("containerInsights"),
				Value: pulumi.String("enabled"),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ECS cluster: %w", err)
	}

	return &ECSClusterOutput{
		Cluster: cluster,
	}, nil
}

func CreateECSFargateService(
	ctx *pulumi.Context,
	cluster *ecs.Cluster,
	imageURI pulumi.StringOutput,
	ecsSubnetID pulumi.Output,
	ecsSecurityGroupID pulumi.StringOutput,
	loadBalancer *lb.ApplicationLoadBalancer,
	rdsEndpoint pulumi.StringOutput,
	rdsPort pulumi.IntOutput,
	rdsUsername pulumi.StringOutput,
	rdsDatabase pulumi.StringOutput,
) (*ECSServiceOutput, error) {
	cfg := config.New(ctx, "")
	environment := cfg.Get("environment")

	// Determine task count based on environment
	desiredCount := 1
	if environment == "prod" {
		desiredCount = 2
	} else if environment == "staging" {
		desiredCount = 2
	}

	// Create ECS service with environment variables
	service, err := awsxecs.NewFargateService(ctx, "service", &awsxecs.FargateServiceArgs{
		Cluster:                       cluster.Arn,
		HealthCheckGracePeriodSeconds: pulumi.Int(300),
		NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
			Subnets: pulumi.StringArray{
				ecsSubnetID.ApplyT(func(id interface{}) string {
					return id.(string)
				}).(pulumi.StringOutput),
			},
			SecurityGroups: pulumi.StringArray{
				ecsSecurityGroupID,
			},
			AssignPublicIp: pulumi.Bool(false),
		},
		DesiredCount: pulumi.Int(desiredCount),
		TaskDefinitionArgs: &awsxecs.FargateServiceTaskDefinitionArgs{
			Container: &awsxecs.TaskDefinitionContainerDefinitionArgs{
				Image:  imageURI,
				Cpu:    pulumi.Int(256),
				Memory: pulumi.Int(512),
				PortMappings: ecsx.TaskDefinitionPortMappingArray{
					ecsx.TaskDefinitionPortMappingArgs{
						ContainerPort: pulumi.Int(8080),
						TargetGroup:   loadBalancer.DefaultTargetGroup,
					},
				},
				// Environment variables for database connection
				Environment: ecsx.TaskDefinitionKeyValuePairArray{
					&ecsx.TaskDefinitionKeyValuePairArgs{
						Name:  pulumi.String("ENVIRONMENT"),
						Value: pulumi.String(environment),
					},
					&ecsx.TaskDefinitionKeyValuePairArgs{
						Name:  pulumi.String("DB_HOST"),
						Value: rdsEndpoint,
					},
					&ecsx.TaskDefinitionKeyValuePairArgs{
						Name: pulumi.String("DB_PORT"),
						Value: rdsPort.ApplyT(func(p interface{}) string {
							return fmt.Sprintf("%v", p)
						}).(pulumi.StringOutput),
					},
					&ecsx.TaskDefinitionKeyValuePairArgs{
						Name:  pulumi.String("DB_USER"),
						Value: rdsUsername,
					},
					&ecsx.TaskDefinitionKeyValuePairArgs{
						Name:  pulumi.String("DB_NAME"),
						Value: rdsDatabase,
					},
					&ecsx.TaskDefinitionKeyValuePairArgs{
						Name:  pulumi.String("PORT"),
						Value: pulumi.String("8080"),
					},
				},
				// CloudWatch logging configuration
				LogConfiguration: &ecsx.TaskDefinitionLogConfigurationArgs{
					LogDriver: pulumi.String("awslogs"),
					Options: pulumi.StringMap{
						"awslogs-group":         pulumi.Sprintf("/ecs/ecom-%s", environment),
						"awslogs-region":        pulumi.String("sa-east-1"),
						"awslogs-stream-prefix": pulumi.String("ecs"),
					},
				},
			},
		},
		Tags: pulumi.StringMap{
			"Name":        pulumi.String("ecom-service"),
			"Environment": pulumi.String(environment),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ECS Fargate service: %w", err)
	}

	ctx.Log.Info(fmt.Sprintf("ECS Service created with %d desired tasks for %s environment", desiredCount, environment), nil)

	return &ECSServiceOutput{
		Service: service,
	}, nil
}
