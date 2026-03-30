# Regional Migration - US-EAST-1 to SA-EAST-1

**Current Status:**
- Infrastructure exists in: **us-east-1** (Virginia - USA)
- Target region: **sa-east-1** (São Paulo - Brazil)
- Log group: ✅ Created in both regions

## Current State

### US-EAST-1 (Existing)
```
VPC:                  vpc-0be9581cd20d92537
Subnets:              6 subnets (3 public, 3 private)
RDS:                  ecom-databasea228fa2.cu984u0uqulk.us-east-1.rds.amazonaws.com
Load Balancer:        lb-872b917-361446910.us-east-1.elb.amazonaws.com
ECS Cluster:          cluster-205c873
CloudWatch Logs:      ecom-api-logs (exists)
```

### SA-EAST-1 (New Target)
```
VPC:                  (none yet)
CloudWatch Logs:      ecom-api-logs (pre-created)
```

## Migration Path

### Option A: Clean Migration (Recommended)
1. **Destroy us-east-1 stack** - Clean up all resources
2. **Deploy fresh to sa-east-1** - New infrastructure in São Paulo
3. Benefits: Clean slate, no orphaned resources

### Option B: Keep Both (Testing)
1. **Keep us-east-1 running** - Keep for reference/testing
2. **Deploy new to sa-east-1** - New production in São Paulo
3. Benefits: Can compare, easy rollback

## Steps for Option A (Recommended)

### Step 1: Destroy us-east-1 Infrastructure
```bash
cd infra/pulumi
pulumi stack select dev
pulumi destroy --yes --region us-east-1
```

### Step 2: Push Updated Code (sa-east-1)
```bash
git push origin main
```

This triggers GitHub Actions which will:
1. Run tests
2. Build Docker image
3. Deploy to ECR in sa-east-1
4. Deploy infrastructure to sa-east-1 via Pulumi

### Step 3: Get New Load Balancer DNS
```bash
cd infra/pulumi
pulumi stack output loadBalancerDns
# Will return: lb-XXX-XXX.sa-east-1.elb.amazonaws.com
```

### Step 4: Test API
```bash
http http://lb-XXX-XXX.sa-east-1.elb.amazonaws.com/health
```

## Comparison: us-east-1 vs sa-east-1

| Aspect | us-east-1 | sa-east-1 |
|--------|-----------|-----------|
| Region Name | N. Virginia | São Paulo |
| Latency | Higher for Brazil | Lower (local) |
| DNS Endpoint | elb.us-east-1.amazonaws.com | elb.sa-east-1.amazonaws.com |
| ECR Repository | ecr.us-east-1.amazonaws.com | ecr.sa-east-1.amazonaws.com |
| RDS Endpoint | rds.us-east-1.amazonaws.com | rds.sa-east-1.amazonaws.com |

## Files Changed for Migration

```
✅ infra/pulumi/resources/ecs.go           - awslogs-region: sa-east-1
✅ .github/workflows/ci-cd.yml             - AWS_REGION: sa-east-1
```

## Next Steps

**Choose your path:**

1. **Recommended: Clean Migration**
   ```bash
   # Destroy old infrastructure
   cd infra/pulumi && pulumi destroy --yes
   
   # Deploy to sa-east-1
   git push origin main
   ```

2. **Or: Keep both running**
   ```bash
   # Just push new code
   git push origin main
   # (GitHub Actions will deploy to sa-east-1)
   ```

---
**Status:** Ready for sa-east-1 deployment ✅
