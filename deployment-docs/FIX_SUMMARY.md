# Multi-AZ Deployment Fix - Summary

**Date:** March 30, 2026  
**Status:** ✅ FIXED - Ready for deployment testing

## Problem
The infrastructure deployment was failing with multi-AZ requirements not being met:
- RDS subnet group creation: "needs 2+ availability zones"
- ALB creation: "needs 2+ subnets across different AZs"

## Root Cause
1. **VPC Configuration:** Duplicate `CreateVPC` function with type errors prevented proper multi-AZ subnet creation
2. **Subnet Distribution:** Even though VPC created 3 AZs, subnets were only being passed one at a time to RDS and ALB

## Fixes Applied

### 1. Fixed vpc.go (Commit: 1473e25)
**Issue:** Duplicate `CreateVPC` function declaration with type errors
- Removed duplicate function definition
- Fixed `AvailabilityZoneNames` parameter: `pulumi.StringArray` → `[]string`
- Removed invalid `EnableDns` and `EnableDnsHostnames` fields (not in VpcArgs)
- Kept explicit 3-AZ configuration: `us-east-1a`, `us-east-1b`, `us-east-1c`
- Configured NAT gateway with `OnePerAz` strategy

**Result:** Code now builds successfully ✓

### 2. Fixed multi-AZ subnet distribution (Commit: 9d92fb8)
**Issue:** Only single subnets were being passed to ALB and RDS
- Updated `main.go`:
  - Pass `vpcOutput.PublicSubnets` (all public subnets) to LoadBalancer
  - Pass `vpcOutput.PrivateSubnets` (all private subnets) to RDS
  - Removed single subnet extractions (`.Index(0)`, `.Index(1)`)

- Updated `loadbalancer.go`:
  - Function now accepts `pulumi.StringArrayOutput` instead of single subnet
  - Removed manual subnet ID conversion logic

- Updated `rds.go`:
  - Function now accepts `pulumi.StringArrayOutput` instead of single subnet
  - RDS SubnetGroup now receives all private subnets (multi-AZ compliant)

**Result:** Multi-AZ requirements satisfied ✓

## Infrastructure Architecture (Post-Fix)

```
VPC (10.0.0.0/16)
├── 3 Availability Zones (us-east-1a, b, c)
│
├── Public Subnets (1 per AZ × 3)
│   └── Application Load Balancer (spans all 3 AZs)
│       ├── Target Group for ECS Fargate
│       └── Health checks on port 8080/health
│
├── Private Subnets (1 per AZ × 3)
│   ├── RDS PostgreSQL DB Subnet Group (multi-AZ)
│   │   └── Primary in 1 AZ + standby replicas in other AZs
│   ├── ECS Fargate Tasks
│   └── NAT Gateways (1 per AZ)
│
└── Security Groups
    ├── ALB SG: Allow 80 (HTTP) from internet
    ├── ECS SG: Allow 8080 from ALB, allow all to RDS
    └── RDS SG: Allow 5432 from ECS SG
```

## Pre-Deployment Checklist

- [x] VPC code builds without errors
- [x] Multi-AZ subnet configuration correct
- [x] All resources have access to 3 AZs
- [x] RDS subnet group will span 3 AZs
- [x] ALB will span 3 AZs
- [x] ECS Fargate can scale across AZs
- [x] NAT gateway strategy configured (OnePerAz)
- [x] GitHub Actions workflow ready
- [x] AWS credentials configured in GitHub Secrets
- [ ] **NEXT:** Test deployment via `pulumi up`

## Next Steps

1. **Verify AWS credentials and permissions:**
   ```bash
   aws sts get-caller-identity
   ```

2. **Test deployment locally (optional):**
   ```bash
   cd infra/pulumi
   pulumi stack select dev
   pulumi preview
   pulumi up
   ```

3. **Push changes to trigger GitHub Actions:**
   ```bash
   git push origin main
   ```
   This will trigger:
   - Test Stage: Run Go tests against PostgreSQL
   - Build Stage: Build Docker image and push to ECR
   - Deploy Stage: Deploy infrastructure with Pulumi

4. **Monitor deployment:**
   - Check GitHub Actions workflow
   - Verify resources in AWS Console:
     - VPC with 3 AZs
     - RDS instance with standby replicas
     - ALB across 3 AZs
     - ECS Fargate service running

## Key Files Modified

| File | Change | Purpose |
|------|--------|---------|
| `infra/pulumi/resources/vpc.go` | Removed duplicate, fixed types | Create 3-AZ VPC correctly |
| `infra/pulumi/main.go` | Pass all subnets to resources | Distribute across AZs |
| `infra/pulumi/resources/loadbalancer.go` | Accept StringArrayOutput | ALB spans all public subnets |
| `infra/pulumi/resources/rds.go` | Accept StringArrayOutput | RDS uses all private subnets |

## Commits

```
9d92fb8 fix: Pass all subnets to ALB and RDS for true multi-AZ deployment
1473e25 fix: Remove duplicate CreateVPC function and fix Pulumi type errors
```

## Validation Commands

```bash
# Build Pulumi code
cd infra/pulumi
go build ./...

# Verify Go version
go version

# Check Docker file
docker build --no-cache -t ecom-api:test .

# Inspect VPC configuration (after deployment)
aws ec2 describe-subnets --filters "Name=vpc-id,Values=<vpc-id>"
aws rds describe-db-instances --db-instance-identifier ecom-database
aws elbv2 describe-load-balancers --names lb
```

---
**Status:** Ready for production deployment testing ✅
