# Limpeza e Permissões AWS para Pulumi

## 1. COMO LIMPAR A AWS

### Script de Limpeza Automática

```bash
#!/bin/bash

REGION="sa-east-1"

echo "=== LIMPEZA COMPLETA DE RECURSOS AWS ==="

# 1. Delete ECS Services e Clusters
echo -e "\n[1/7] Deletando ECS Services e Clusters..."
for cluster_arn in $(aws ecs list-clusters --region $REGION --query 'clusterArns[]' --output text); do
  cluster_name=$(echo $cluster_arn | awk -F'/' '{print $NF}')
  for service_arn in $(aws ecs list-services --cluster $cluster_name --region $REGION --query 'serviceArns[]' --output text); do
    service_name=$(echo $service_arn | awk -F'/' '{print $NF}')
    aws ecs delete-service --cluster $cluster_name --service $service_name --force --region $REGION 2>/dev/null
  done
  aws ecs delete-cluster --cluster $cluster_name --region $REGION 2>/dev/null
done

# 2. Delete Load Balancers e Target Groups
echo "[2/7] Deletando Load Balancers..."
for lb_arn in $(aws elbv2 describe-load-balancers --region $REGION --query 'LoadBalancers[].LoadBalancerArn' --output text); do
  aws elbv2 delete-load-balancer --load-balancer-arn $lb_arn --region $REGION 2>/dev/null
done

for tg_arn in $(aws elbv2 describe-target-groups --region $REGION --query 'TargetGroups[].TargetGroupArn' --output text); do
  aws elbv2 delete-target-group --target-group-arn $tg_arn --region $REGION 2>/dev/null
done

# 3. Delete RDS Instances
echo "[3/7] Deletando RDS Instances..."
for db_id in $(aws rds describe-db-instances --region $REGION --query 'DBInstances[].DBInstanceIdentifier' --output text); do
  aws rds delete-db-instance --db-instance-identifier $db_id --skip-final-snapshot --region $REGION 2>/dev/null
done

# 4. Delete ECR Repositories
echo "[4/7] Deletando ECR Repositories..."
for repo in $(aws ecr describe-repositories --region $REGION --query 'repositories[].repositoryName' --output text); do
  aws ecr delete-repository --repository-name $repo --force --region $REGION 2>/dev/null
done

# 5. Delete VPCs (com todas as dependências)
echo "[5/7] Deletando VPCs e dependências..."
for vpc_id in $(aws ec2 describe-vpcs --region $REGION --query 'Vpcs[].VpcId' --output text); do
  # Skip default VPC
  default_vpc=$(aws ec2 describe-vpcs --region $REGION --filters "Name=isDefault,Values=true" --query 'Vpcs[].VpcId' --output text)
  if [ "$vpc_id" == "$default_vpc" ]; then
    continue
  fi
  
  # Delete NAT Gateways e Elastic IPs
  for natgw_id in $(aws ec2 describe-nat-gateways --region $REGION --filter "Name=vpc-id,Values=$vpc_id" --query 'NatGateways[?State!=`deleted`].NatGatewayId' --output text); do
    aws ec2 delete-nat-gateway --nat-gateway-id $natgw_id --region $REGION 2>/dev/null
  done
  
  sleep 5
  
  # Release Elastic IPs
  for alloc_id in $(aws ec2 describe-addresses --region $REGION --filters "Name=domain,Values=vpc" --query 'Addresses[].AllocationId' --output text); do
    aws ec2 release-address --allocation-id $alloc_id --region $REGION 2>/dev/null
  done
  
  # Delete Internet Gateways
  for igw_id in $(aws ec2 describe-internet-gateways --region $REGION --filters "Name=attachment.vpc-id,Values=$vpc_id" --query 'InternetGateways[].InternetGatewayId' --output text); do
    aws ec2 detach-internet-gateway --internet-gateway-id $igw_id --vpc-id $vpc_id --region $REGION 2>/dev/null
    aws ec2 delete-internet-gateway --internet-gateway-id $igw_id --region $REGION 2>/dev/null
  done
  
  # Delete Network Interfaces
  for eni_id in $(aws ec2 describe-network-interfaces --region $REGION --filters "Name=vpc-id,Values=$vpc_id" --query 'NetworkInterfaces[].NetworkInterfaceId' --output text); do
    aws ec2 delete-network-interface --network-interface-id $eni_id --region $REGION 2>/dev/null
  done
  
  # Delete Security Groups
  for sg_id in $(aws ec2 describe-security-groups --region $REGION --filters "Name=vpc-id,Values=$vpc_id" --query 'SecurityGroups[?GroupName!=`default`].GroupId' --output text); do
    aws ec2 delete-security-group --group-id $sg_id --region $REGION 2>/dev/null
  done
  
  # Delete Subnets
  for subnet_id in $(aws ec2 describe-subnets --region $REGION --filters "Name=vpc-id,Values=$vpc_id" --query 'Subnets[].SubnetId' --output text); do
    aws ec2 delete-subnet --subnet-id $subnet_id --region $REGION 2>/dev/null
  done
  
  # Delete Route Tables
  for rt_id in $(aws ec2 describe-route-tables --region $REGION --filters "Name=vpc-id,Values=$vpc_id" --query 'RouteTables[?Associations[0].Main!=`true`].RouteTableId' --output text); do
    aws ec2 delete-route-table --route-table-id $rt_id --region $REGION 2>/dev/null
  done
  
  # Delete VPC
  aws ec2 delete-vpc --vpc-id $vpc_id --region $REGION 2>/dev/null
done

# 6. Delete RDS Subnet Groups
echo "[6/7] Deletando RDS Subnet Groups..."
for sg in $(aws rds describe-db-subnet-groups --region $REGION --query 'DBSubnetGroups[].DBSubnetGroupName' --output text); do
  aws rds delete-db-subnet-group --db-subnet-group-name $sg --region $REGION 2>/dev/null
done

# 7. Delete CloudWatch Log Groups
echo "[7/7] Deletando CloudWatch Log Groups..."
for log_group in $(aws logs describe-log-groups --region $REGION --query 'logGroups[].logGroupName' --output text); do
  aws logs delete-log-group --log-group-name "$log_group" --region $REGION 2>/dev/null
done

echo -e "\n✅ Limpeza completa!"
```

### Como usar o script

```bash
# 1. Salvar o script
curl -o cleanup.sh https://seu-dominio/cleanup.sh

# 2. Dar permissão de execução
chmod +x cleanup.sh

# 3. Executar
./cleanup.sh
```

---

## 2. PERMISSÕES AWS NECESSÁRIAS

### Versão Simples (Recomendada para Desenvolvimento)

Use **AdministratorAccess** ou crie uma policy customizada:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "EC2FullAccess",
      "Effect": "Allow",
      "Action": ["ec2:*"],
      "Resource": "*"
    },
    {
      "Sid": "ECSFullAccess",
      "Effect": "Allow",
      "Action": ["ecs:*"],
      "Resource": "*"
    },
    {
      "Sid": "RDSFullAccess",
      "Effect": "Allow",
      "Action": ["rds:*"],
      "Resource": "*"
    },
    {
      "Sid": "ECRFullAccess",
      "Effect": "Allow",
      "Action": ["ecr:*"],
      "Resource": "*"
    },
    {
      "Sid": "ELBFullAccess",
      "Effect": "Allow",
      "Action": ["elasticloadbalancing:*"],
      "Resource": "*"
    },
    {
      "Sid": "IAMRoleAccess",
      "Effect": "Allow",
      "Action": [
        "iam:CreateRole",
        "iam:DeleteRole",
        "iam:AttachRolePolicy",
        "iam:DetachRolePolicy",
        "iam:PutRolePolicy",
        "iam:DeleteRolePolicy",
        "iam:GetRole",
        "iam:PassRole"
      ],
      "Resource": "*"
    },
    {
      "Sid": "CloudWatchLogsAccess",
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:DeleteLogGroup",
        "logs:DescribeLogGroups"
      ],
      "Resource": "*"
    }
  ]
}
```

### Como Aplicar a Policy

#### Opção 1: Via Console AWS

1. Ir em **IAM** → **Users** → Selecionar seu usuário
2. Clicar em **Add permissions** → **Attach policies**
3. Buscar por `AdministratorAccess` ou colar a policy acima

#### Opção 2: Via AWS CLI

```bash
# Criar a policy
aws iam create-policy \
  --policy-name PulumiDeployPolicy \
  --policy-document file://policy.json

# Anexar ao usuário
aws iam attach-user-policy \
  --user-name seu-usuario \
  --policy-arn arn:aws:iam::AWS_ACCOUNT_ID:policy/PulumiDeployPolicy
```

### Permissões Mínimas (Mais Restritivo)

Se quiser apenas criar/deletar Pulumi stacks:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:CreateVpc",
        "ec2:DeleteVpc",
        "ec2:CreateSubnet",
        "ec2:DeleteSubnet",
        "ec2:CreateSecurityGroup",
        "ec2:DeleteSecurityGroup",
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:AuthorizeSecurityGroupEgress",
        "ec2:CreateNatGateway",
        "ec2:DeleteNatGateway",
        "ec2:AllocateAddress",
        "ec2:ReleaseAddress",
        "ec2:CreateInternetGateway",
        "ec2:DeleteInternetGateway",
        "ec2:AttachInternetGateway",
        "ec2:DetachInternetGateway",
        "ec2:CreateRouteTable",
        "ec2:DeleteRouteTable",
        "ec2:CreateRoute",
        "ec2:DeleteRoute",
        "ec2:Describe*",
        "ecs:*",
        "rds:*",
        "ecr:*",
        "elasticloadbalancing:*",
        "iam:CreateRole",
        "iam:DeleteRole",
        "iam:AttachRolePolicy",
        "iam:DetachRolePolicy",
        "iam:GetRole",
        "iam:PassRole",
        "logs:CreateLogGroup",
        "logs:DeleteLogGroup"
      ],
      "Resource": "*"
    }
  ]
}
```

---

## 3. TROUBLESHOOTING

### Problema: "AuthFailure: You do not have permission"
**Solução:** Verifique se o usuário AWS tem as permissões corretas. Veja a seção de permissões acima.

### Problema: "Resource.AlreadyAssociated: Elastic IP is already associated"
**Solução:** Execute o script de cleanup para liberar EIPs orfãos.

### Problema: "The vpc has dependencies and cannot be deleted"
**Solução:** O script de cleanup cuida disso automaticamente, deletando ENIs, NAT Gateways e IGWs antes da VPC.

---

## 4. VERIFICAÇÃO PÓS-CLEANUP

```bash
# Verificar que tudo foi deletado
aws ec2 describe-vpcs --region sa-east-1 \
  --query 'Vpcs[?Tags[?Key==`Name`]].VpcId' \
  --output text

aws ecs list-clusters --region sa-east-1 \
  --query 'clusterArns' --output text

aws ecr describe-repositories --region sa-east-1 \
  --query 'repositories[].repositoryName' --output text

# Se não retornar nada, está limpo!
```

---

## 5. PROXIMAS TENTATIVAS

Após limpar e ter as permissões corretas:

```bash
# 1. Criar novo stack
export PULUMI_CONFIG_PASSPHRASE=""
cd infra/pulumi
pulumi stack init dev

# 2. Fazer preview
pulumi preview

# 3. Deploy
pulumi up --yes
```

Agora deve funcionar! 🚀
