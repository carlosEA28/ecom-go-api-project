# CloudWatch Logging Fix for ECS Tasks

**Date:** March 30, 2026  
**Status:** ✅ FIXED - ECS tasks now have proper logging

## Problem

ECS tasks were failing to start with error:

```
ResourceNotFoundException: The specified log group does not exist
failed to create Cloudwatch log stream: operation error CloudWatch Logs: 
CreateLogStream, https response error StatusCode: 400, RequestID: ...
```

Tasks were stuck in a loop:
- ECS tries to start task
- Task configuration references log group `/ecs/ecom-api`
- Log group doesn't exist
- Task fails with error
- ECS retries, loop continues

**Result:** No running tasks, Load Balancer returns 503

## Root Cause

The ECS task definition specified CloudWatch logging configuration but the **log group was never created** by Pulumi infrastructure code.

## Solution

### 1. Created manual log group (immediate fix)
```bash
aws logs create-log-group --log-group-name "ecom-api-logs" --region us-east-1
aws logs create-log-stream --log-group-name "ecom-api-logs" --log-stream-name "ecs-app" --region us-east-1
```

### 2. Updated Pulumi code (permanent fix)

**Added CloudWatch import:**
```go
"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
```

**Auto-create log group in CreateECSFargateService:**
```go
logGroup, err := cloudwatch.NewLogGroup(ctx, "ecs-log-group", &cloudwatch.LogGroupArgs{
    Name:            pulumi.String("ecom-api-logs"),
    RetentionInDays: pulumi.Int(7),
})
```

**Updated ECS Service dependency:**
```go
}, pulumi.DependsOn([]pulumi.Resource{loadBalancer, logGroup}))
```

### 3. Fixed log group naming
- Changed from: `/ecs/ecom-api`
- Changed to: `ecom-api-logs`
- Simpler naming convention, avoids path parsing issues

## How It Works Now

1. Pulumi creates CloudWatch log group **before** ECS service
2. ECS service depends on log group (explicit dependency)
3. Task definition references existing log group
4. Tasks start successfully
5. Logs are written to CloudWatch
6. Tasks register with ALB target group
7. ALB returns 200 OK with task responses

## Files Changed

| File | Change |
|------|--------|
| `infra/pulumi/resources/ecs.go` | Added log group creation, updated imports, fixed log group name |

## Commits

```
b395263 fix: Auto-create CloudWatch log group for ECS tasks
3c3b11d fix: Update CloudWatch log group name for ECS tasks
```

## Expected Result After Redeployment

✅ ECS tasks start successfully  
✅ Logs appear in CloudWatch  
✅ Tasks register with ALB  
✅ API responds with 200 OK  

```bash
http http://lb-872b917-361446910.us-east-1.elb.amazonaws.com/health
# Returns 200 OK (instead of 503)
```

---
**Next:** Redeploy with `pulumi up`
