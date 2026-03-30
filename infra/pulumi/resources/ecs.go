package resources

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecs"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ECSClusterOutput struct {
	Cluster *ecs.Cluster
}

type ECSServiceOutput struct {
	Service *ecs.Service
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
	// NOTE: Update awslogs-region below if deploying to a different AWS region
	// Supported values: us-east-1, sa-east-1, eu-west-1, ap-southeast-1, etc.
	awsRegion := "us-east-1"

	// Log group creation is handled externally or via Pulumi import
	// This avoids conflicts if the log group was created manually
	// The ECS task definition references "ecom-api-logs" directly

	// Create execution role for ECS tasks
	executionRole, err := iam.NewRole(ctx, "ecs-execution-role", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
					"Service": "ecs-tasks.amazonaws.com"
				}
			}]
		}`),
	})
	if err != nil {
		return nil, err
	}

	// Attach execution role policy
	_, err = iam.NewRolePolicyAttachment(ctx, "ecs-execution-role-policy", &iam.RolePolicyAttachmentArgs{
		Role:      executionRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"),
	})
	if err != nil {
		return nil, err
	}

	// Create task definition
	taskDef, err := ecs.NewTaskDefinition(ctx, "app-task", &ecs.TaskDefinitionArgs{
		Family:                  pulumi.String("ecom-api"),
		NetworkMode:             pulumi.String("awsvpc"),
		RequiresCompatibilities: pulumi.StringArray{pulumi.String("FARGATE")},
		Cpu:                     pulumi.String("256"),
		Memory:                  pulumi.String("512"),
		ExecutionRoleArn:        executionRole.Arn,
		ContainerDefinitions: imageURI.ApplyT(func(uri interface{}) string {
			containers := []map[string]interface{}{
				{
					"name":      "ecom-api",
					"image":     uri.(string),
					"essential": true,
					"portMappings": []map[string]interface{}{
						{
							"containerPort": 8080,
							"hostPort":      8080,
							"protocol":      "tcp",
						},
					},
					"logConfiguration": map[string]interface{}{
						"logDriver": "awslogs",
						"options": map[string]interface{}{
							"awslogs-group":         "ecom-api-logs",
							"awslogs-region":        awsRegion,
							"awslogs-stream-prefix": "ecs",
						},
					},
				},
			}
			jsonBytes, _ := json.Marshal(containers)
			return string(jsonBytes)
		}).(pulumi.StringOutput),
	})
	if err != nil {
		return nil, err
	}

	// Create ECS service
	service, err := ecs.NewService(ctx, "app-service", &ecs.ServiceArgs{
		Cluster:        cluster.Arn,
		TaskDefinition: taskDef.Arn,
		DesiredCount:   pulumi.Int(1),
		LaunchType:     pulumi.String("FARGATE"),
		NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
			Subnets:        ecsSubnetIDs,
			SecurityGroups: pulumi.StringArray{ecsSecurityGroupID},
			AssignPublicIp: pulumi.Bool(false),
		},
		LoadBalancers: ecs.ServiceLoadBalancerArray{
			&ecs.ServiceLoadBalancerArgs{
				TargetGroupArn: loadBalancer.DefaultTargetGroup.Arn(),
				ContainerName:  pulumi.String("ecom-api"),
				ContainerPort:  pulumi.Int(8080),
			},
		},
		HealthCheckGracePeriodSeconds: pulumi.Int(300),
	}, pulumi.DependsOn([]pulumi.Resource{loadBalancer}))
	if err != nil {
		return nil, fmt.Errorf("error creating ECS service: %w", err)
	}

	return &ECSServiceOutput{
		Service: service,
	}, nil
}
