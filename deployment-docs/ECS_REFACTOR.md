# ECS Implementation Fix - Raw AWS Resources

**Date:** March 30, 2026  
**Status:** ✅ REFACTORED - Ready for clean deployment

## Problem

The previous deployment failed because:
```
ClientException: When networkMode=awsvpc, the host ports and container 
ports in port mappings must match.
```

Despite adding `HostPort: 8080` to the awsx wrapper configuration, the error persisted because the wrapper was creating task definitions that didn't match the expected structure.

## Root Cause Analysis

The `awsx.FargateService` high-level wrapper was:
1. Hiding the actual task definition structure
2. Creating implicit defaults that conflicted with `awsvpc` mode requirements
3. Providing no clear way to explicitly set host port matching

## Solution: Raw AWS Resources

**Switched from:**
- `awsx.FargateService` (high-level wrapper with limited control)

**Switched to:**
- `aws.ecs.TaskDefinition` (explicit task definition)
- `aws.ecs.Service` (explicit service configuration)
- `aws.iam.Role` (explicit IAM execution role)

## Implementation Details

### Task Definition (Raw JSON)
```go
ContainerDefinitions: imageURI.ApplyT(func(uri interface{}) string {
    containers := []map[string]interface{}{
        {
            "name":      "ecom-api",
            "image":     uri,
            "essential": true,
            "portMappings": []map[string]interface{}{
                {
                    "containerPort": 8080,
                    "hostPort":      8080,  // ✓ Explicit matching
                    "protocol":      "tcp",
                },
            },
            "logConfiguration": {
                "logDriver": "awslogs",
                "options": {
                    "awslogs-group":         "/ecs/ecom-api",
                    "awslogs-region":        "us-east-1",
                    "awslogs-stream-prefix": "ecs",
                },
            },
        },
    }
    jsonBytes, _ := json.Marshal(containers)
    return string(jsonBytes)
}).(pulumi.StringOutput)
```

### Service Configuration
- **NetworkMode:** `awsvpc` (required for Fargate)
- **LaunchType:** `FARGATE`
- **Subnets:** All private subnets (multi-AZ scaling)
- **Security Groups:** ECS security group
- **Target Group:** Connected to ALB
- **Desired Count:** 1 (auto-scaling ready)

### IAM Execution Role
- Automatically created role for ECS task execution
- Attached policy: `AmazonECSTaskExecutionRolePolicy`
- Allows ECS to pull images and log to CloudWatch

## Files Changed

| File | Change |
|------|--------|
| `infra/pulumi/resources/ecs.go` | Completely rewritten with raw AWS resources |
| `infra/pulumi/main.go` | Updated export to use `Service.Arn` (not `Service.Service.Arn()`) |

## Pre-Deployment Cleanup

Before redeploying with the new ECS code, **destroy the previous failed stack**:

```bash
cd infra/pulumi
pulumi stack select dev
pulumi destroy --yes
```

This will remove:
- Failed ECS task definition
- IAM roles (old ones)
- All other resources created so far

Then deploy fresh:
```bash
pulumi up
```

## Advantages of Raw Resources

✅ **Explicit Control:** Every aspect of task definition is clear  
✅ **No Wrapper Abstraction:** Direct mapping to AWS API  
✅ **Port Mapping Clarity:** Host and container ports explicitly defined  
✅ **Error Debugging:** Stack traces directly reference task definition issues  
✅ **Flexibility:** Easy to add more container configurations later  
✅ **Logging:** Built-in CloudWatch logging configuration  

## Expected Behavior After Fix

1. **Task Definition Creation:** Will succeed (ports match)
2. **ECS Service Launch:** Will register tasks with ALB
3. **Container Health:** Container will listen on port 8080
4. **Load Balancer:** Will forward traffic to container on 8080
5. **Multi-AZ:** Tasks can scale across all 3 availability zones

## Commit

```
81b5a32 refactor: Switch ECS to use raw AWS resources instead of awsx wrapper
```

## Next Steps

1. Destroy previous failed stack
2. Deploy with new ECS configuration
3. Verify all resources created successfully
4. Test API endpoints through load balancer

---
**Status:** Ready for clean deployment ✅
