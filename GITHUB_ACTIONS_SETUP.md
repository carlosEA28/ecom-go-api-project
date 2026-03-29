# GitHub Actions Setup - Complete Guide

## 🎯 Overview

Your ecom-go-api project is configured with automated CI/CD pipelines that:

1. **Test** - Runs unit tests on every push/PR
2. **Build** - Compiles the Go application
3. **Deploy** (main branch only) - Deploys infrastructure to AWS and runs migrations

All infrastructure uses:
- **Region**: `sa-east-1` (São Paulo, Brazil)
- **Backend**: Local Pulumi state (stored in `.pulumi` directory)
- **Database**: PostgreSQL 15 on RDS
- **Container**: Docker image in ECR
- **Compute**: ECS Fargate service

## ⚙️ Required GitHub Secrets

Add these secrets to your repository at: **Settings → Secrets and variables → Actions**

### 1. AWS Access Key ID
```
Name: AWS_ACCESS_KEY_ID
Value: <your-aws-access-key-id>
```

### 2. AWS Secret Access Key
```
Name: AWS_SECRET_ACCESS_KEY
Value: <your-aws-secret-access-key>
```

### 3. RDS Master Password
```
Name: RDS_PASSWORD
Value: <strong-postgres-password>
```
⚠️ **Important**: Use a strong password (min 8 chars, mix of upper/lower/numbers/symbols)

### 4. Pulumi Config Passphrase
```
Name: PULUMI_CONFIG_PASSPHRASE
Value: 08h0gEZTuIAwy7jxIp42NcyovBQqixiY
```
⚠️ **CRITICAL**: This exact passphrase is required to decrypt Pulumi secrets during deployment

## 🔐 How to Generate AWS Credentials

### Step 1: Create IAM User (if you don't have one)

1. Go to [AWS IAM Console](https://console.aws.amazon.com/iam/)
2. Click **Users** in the left menu
3. Click **Create user**
4. Enter username (e.g., `ecom-ci-cd-user`)
5. Click **Next**
6. Attach the policy below as "Inline policy"

### Step 2: Create Access Keys

1. Select your IAM user
2. Go to **Security credentials** tab
3. Scroll to **Access keys** section
4. Click **Create access key**
5. Choose **Application running outside AWS**
6. Click **Next**
7. Copy the **Access Key ID** and **Secret Access Key**
8. ⚠️ Save these safely - you won't see the secret again!

### Step 3: Attach Minimal IAM Policy

Replace `{ACCOUNT_ID}` with your AWS Account ID (found in top-right of AWS Console):

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "EC2Permissions",
      "Effect": "Allow",
      "Action": [
        "ec2:*"
      ],
      "Resource": "*"
    },
    {
      "Sid": "RDSPermissions",
      "Effect": "Allow",
      "Action": [
        "rds:*"
      ],
      "Resource": "*"
    },
    {
      "Sid": "ECRPermissions",
      "Effect": "Allow",
      "Action": [
        "ecr:*"
      ],
      "Resource": "*"
    },
    {
      "Sid": "ECSPermissions",
      "Effect": "Allow",
      "Action": [
        "ecs:*"
      ],
      "Resource": "*"
    },
    {
      "Sid": "ELBPermissions",
      "Effect": "Allow",
      "Action": [
        "elasticloadbalancing:*"
      ],
      "Resource": "*"
    },
    {
      "Sid": "IAMPermissions",
      "Effect": "Allow",
      "Action": [
        "iam:GetRole",
        "iam:PassRole",
        "iam:CreateRole",
        "iam:PutRolePolicy",
        "iam:DeleteRolePolicy",
        "iam:DeleteRole",
        "iam:GetRolePolicy"
      ],
      "Resource": "*"
    },
    {
      "Sid": "CloudWatchPermissions",
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "logs:DescribeLogGroups"
      ],
      "Resource": "*"
    }
  ]
}
```

## 📋 Step-by-Step Setup Instructions

### 1. Add GitHub Secrets

**Via GitHub Web Interface:**

1. Go to your repository: https://github.com/carlosEA28/ecom-go-api-project
2. Click **Settings** tab
3. Click **Secrets and variables** → **Actions** in left menu
4. Click **New repository secret** for each secret:

```
1. AWS_ACCESS_KEY_ID = <your-key>
2. AWS_SECRET_ACCESS_KEY = <your-secret>
3. RDS_PASSWORD = <your-strong-password>
4. PULUMI_CONFIG_PASSPHRASE = 08h0gEZTuIAwy7jxIp42NcyovBQqixiY
```

### 2. Push Changes to GitHub

```bash
# View what's ready to commit
git status

# Push to origin (your fork)
git push origin main
```

### 3. Monitor Deployment

1. Go to **Actions** tab in your repository
2. Click the latest workflow run
3. Watch the CI/CD pipeline:
   - ✅ **Test job** - Runs unit tests
   - ✅ **Build job** - Compiles application
   - ✅ **Deploy job** (main only) - Deploys to AWS

## 🚀 What Happens During Deployment

### Deploy Job Flow

```
1. Checkout code
   ↓
2. Configure AWS credentials from secrets
   ↓
3. Set up Pulumi with local backend
   ↓
4. Deploy infrastructure with Pulumi:
   - Create VPC with 3 subnets
   - Create security groups
   - Create RDS PostgreSQL 15
   - Create ECR repository
   - Create load balancer
   - Create ECS Fargate service
   ↓
5. Extract RDS endpoint from Pulumi outputs
   ↓
6. Run database migrations with Goose
   ↓
7. Build Docker image
   ↓
8. Push image to ECR
   ↓
9. Update ECS service with new image
   ↓
10. Output deployment URLs
```

### Expected Outputs

After successful deployment, you'll see:

```
=== Deployment Complete ===
Load Balancer URL: ecom-lb-xxxxx.sa-east-1.elb.amazonaws.com
RDS Endpoint: ecom-rds-xxxxx.sa-east-1.rds.amazonaws.com:5432
ECR Repository: xxxxx.dkr.ecr.sa-east-1.amazonaws.com/ecom
```

## 📊 Infrastructure Architecture

```
┌──────────────────────────────────────────────────┐
│              AWS - sa-east-1 (São Paulo)          │
├──────────────────────────────────────────────────┤
│                                                   │
│  ┌─────────────────────────────────────────┐    │
│  │            Internet Gateway              │    │
│  └────────┬────────────────────────────────┘    │
│           │                                      │
│  ┌────────▼──────────────────────────────────┐  │
│  │      Load Balancer (port 80)               │  │
│  │  DNS: ecom-lb-xxxxx.elb.amazonaws.com    │  │
│  └────────┬──────────────────────────────────┘  │
│           │                                      │
│  ┌────────▼──────────────────────────────────┐  │
│  │         Public Subnet                      │  │
│  │  CIDR: 10.0.0.0/26                        │  │
│  └────────┬──────────────────────────────────┘  │
│           │                                      │
│           │ (Route to private subnets)          │
│           │                                      │
│  ┌────────▼──────────────────────────────────┐  │
│  │         Private Subnet (ECS)               │  │
│  │  CIDR: 10.0.0.64/26                       │  │
│  │  ┌──────────────────────────────────────┐ │  │
│  │  │  ECS Fargate Service                 │ │  │
│  │  │  - Task CPU: 256                     │ │  │
│  │  │  - Task Memory: 512MB                │ │  │
│  │  │  - Image: ECR ecom:latest            │ │  │
│  │  │  - Port: 8080                        │ │  │
│  │  │  - 1 task running                    │ │  │
│  │  └──────────────────────────────────────┘ │  │
│  └────────────────────────────────────────────┘  │
│           │                                      │
│  ┌────────▼──────────────────────────────────┐  │
│  │         Private Subnet (RDS)               │  │
│  │  CIDR: 10.0.0.128/26                      │  │
│  │  ┌──────────────────────────────────────┐ │  │
│  │  │  RDS PostgreSQL 15                   │ │  │
│  │  │  - Instance: db.t3.micro             │ │  │
│  │  │  - Storage: 20GB                     │ │  │
│  │  │  - Database: ecom                    │ │  │
│  │  │  - Port: 5432                        │ │  │
│  │  └──────────────────────────────────────┘ │  │
│  └────────────────────────────────────────────┘  │
│                                                   │
│  ┌──────────────────────────────────────────┐   │
│  │         ECR Repository                   │   │
│  │  - Repository: ecom                      │   │
│  │  - Image Tag: latest, git-sha            │   │
│  └──────────────────────────────────────────┘   │
│                                                   │
└──────────────────────────────────────────────────┘
```

## 🔍 Troubleshooting

### Issue: "AWS credentials not configured"

**Cause**: Secrets are not properly added to the repository

**Fix**:
1. Go to **Settings** → **Secrets and variables** → **Actions**
2. Verify all 4 secrets are present:
   - ✅ AWS_ACCESS_KEY_ID
   - ✅ AWS_SECRET_ACCESS_KEY
   - ✅ RDS_PASSWORD
   - ✅ PULUMI_CONFIG_PASSPHRASE

### Issue: "Pulumi decryption failed"

**Cause**: Wrong or missing PULUMI_CONFIG_PASSPHRASE

**Fix**:
1. Update the secret to: `08h0gEZTuIAwy7jxIp42NcyovBQqixiY`
2. Re-run the workflow

### Issue: "RDS connection timeout"

**Cause**: RDS takes time to become available (usually 2-5 minutes)

**Fix**:
- Workflow includes 30-second retry loop (300 seconds total)
- Check RDS status in AWS Console
- Workflow should eventually succeed

### Issue: "ECR repository not found"

**Cause**: ECR wasn't created by Pulumi

**Fix**:
1. Check Pulumi deployment logs
2. Verify AWS credentials have ECR permissions
3. Manually verify in AWS Console: ECR → Repositories

### Issue: "ECS service not updating"

**Cause**: ECS service not properly configured

**Fix**:
1. Check ECS logs in CloudWatch
2. Verify security group allows port 8080
3. Check task definition and service configuration

## 📝 Monitoring and Logs

### GitHub Actions Logs

1. Go to **Actions** tab
2. Click the workflow run you want to inspect
3. Click the job (test, build, or deploy)
4. Click individual step to see detailed logs

### AWS CloudWatch Logs

After deployment, logs appear in CloudWatch:

```
Log Groups:
├── /ecs/ecom-task - ECS task logs
├── /rds/instance/ecom - RDS logs
└── /pulumi/stack/dev - Pulumi logs (if configured)
```

### AWS Console Verification

**Verify Deployment:**

1. **ECR**: https://console.aws.amazon.com/ecr/repositories/ (sa-east-1)
2. **RDS**: https://console.aws.amazon.com/rds/v2/instances/ (sa-east-1)
3. **ECS**: https://console.aws.amazon.com/ecs/v2/clusters/ (sa-east-1)
4. **Load Balancer**: https://console.aws.amazon.com/ec2/v2/home?region=sa-east-1#LoadBalancers

## 🔄 Re-deploying

### Option 1: Push to main branch

```bash
git push origin main
```

Workflow will automatically trigger on push to main.

### Option 2: Manual workflow dispatch

1. Go to **Actions** tab
2. Select **CI/CD Pipeline**
3. Click **Run workflow**
4. Choose branch (main)
5. Click **Run workflow**

## 🛑 Destroying Infrastructure

If you need to destroy all AWS resources:

```bash
# Locally (with AWS credentials configured)
export AWS_ACCESS_KEY_ID=<your-key>
export AWS_SECRET_ACCESS_KEY=<your-secret>
export PULUMI_CONFIG_PASSPHRASE=08h0gEZTuIAwy7jxIp42NcyovBQqixiY

cd infra/pulumi
pulumi stack select dev
pulumi destroy
```

⚠️ **This will delete all AWS resources including RDS data!**

## 📞 Next Steps

1. ✅ Generate AWS credentials
2. ✅ Add all 4 secrets to GitHub
3. ✅ Push code to main branch: `git push origin main`
4. ✅ Watch the CI/CD pipeline in Actions tab
5. ✅ Access your application via Load Balancer URL
6. ✅ Verify database is populated with migrations

## 📚 Useful Links

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Pulumi AWS Documentation](https://www.pulumi.com/docs/reference/pkg/aws/)
- [AWS IAM Best Practices](https://docs.github.com/en/code-security/secret-scanning/protecting-pushes-with-secret-scanning)
- [Our MIGRATIONS.md Guide](./MIGRATIONS.md)
- [Our ENV_VARIABLES.md Reference](./ENV_VARIABLES.md)
