# 🚀 Guia Completo de Deploy na AWS

## 📋 Pré-requisitos

### 1. AWS Account
- [ ] Conta AWS criada
- [ ] Permissões de Admin ou Policy específica (veja ao final)
- [ ] Region definida (default: `us-east-1`)

### 2. Ferramentas Locais
```bash
# AWS CLI
pip install awscli

# Pulumi CLI
# Windows:
choco install pulumi
# Mac:
brew install pulumi
# Linux:
curl -fsSL https://get.pulumi.com | sh

# Docker (para testar build localmente)
# Windows: https://docs.docker.com/desktop/install/windows-install/
# Mac: https://docs.docker.com/desktop/install/mac-install/
# Linux: https://docs.docker.com/engine/install/

# Go 1.25.6 (já tem)
```

### 3. Verificar Instalações
```bash
aws --version
pulumi version
docker --version
go version
```

---

## 🔑 Passo 1: Configurar AWS Credentials

### 1.1 Criar IAM User na AWS
1. Acesse: https://console.aws.amazon.com/iam/
2. Clique em **Users** → **Create user**
3. Nome: `github-deployer` (ou outro)
4. Próximo → Clique em **Attach policies directly**
5. Procure por e selecione:
   - `AdministratorAccess` (para desenvolvimento)
   - OU crie policy customizada (veja ao final)
6. Clique em **Create user**

### 1.2 Gerar Access Key
1. Clique no user criado
2. Abra a aba **Security credentials**
3. Clique em **Create access key**
4. Selecione **Command Line Interface (CLI)**
5. Clique em **Create access key**
6. **COPIE E GUARDE** em local seguro:
   - `AWS_ACCESS_KEY_ID`
   - `AWS_SECRET_ACCESS_KEY`

### 1.3 Configurar AWS CLI Localmente
```bash
aws configure
# AWS Access Key ID [None]: paste_seu_access_key_id
# AWS Secret Access Key [None]: paste_seu_secret_access_key
# Default region name [None]: us-east-1
# Default output format [None]: json
```

### 1.4 Testar Conexão
```bash
aws sts get-caller-identity
```

Deve retornar algo como:
```json
{
    "UserId": "AIDAXXXXX",
    "Account": "123456789012",
    "Arn": "arn:aws:iam::123456789012:user/github-deployer"
}
```

---

## 🏗️ Passo 2: Configurar Pulumi

### 2.1 Criar Pulumi Account
1. Acesse: https://app.pulumi.com/signup
2. Crie uma conta (pode usar GitHub login)
3. Confirme email

### 2.2 Gerar Pulumi Access Token
1. Acesse: https://app.pulumi.com/account/tokens
2. Clique em **Create token**
3. Nome: `github-ci-cd`
4. **COPIE E GUARDE** o token

### 2.3 Fazer Login Localmente
```bash
pulumi login
# Cole o token quando pedido
```

### 2.4 Preparar Pulumi Project
```bash
cd infra/pulumi

# Ver stacks existentes
pulumi stack ls

# Se não tem stack dev, criar:
pulumi stack init dev

# Configurar stack:
pulumi config set aws:region us-east-1
```

### 2.5 Preview das mudanças
```bash
# Ver o que será criado (NÃO cria ainda)
pulumi preview
```

Deve mostrar recursos como:
```
+ aws:vpc:Vpc ...
+ aws:rds:Instance ...
+ aws:ecs:Cluster ...
+ aws:ecr:Repository ...
+ Load Balancer
+ Security Groups
...
```

---

## 📦 Passo 3: Deploy Manual (Teste)

### 3.1 Deploy com Pulumi
```bash
cd infra/pulumi

# Fazer deploy (vai criar TODOS os recursos AWS)
pulumi up
```

Vai mostrar:
- Recursos a serem criados (com custo estimado)
- Pedirá confirmação: `yes` para confirmar

### 3.2 Monitorar Deploy
```bash
# Enquanto deploy está rodando, pode ver logs:
pulumi logs
```

### 3.3 Pegar Outputs
```bash
# Ver todos os outputs
pulumi stack output

# Ver específico:
pulumi stack output loadBalancerDns
# Deve retornar: ecom-nlb-123456-789.elb.us-east-1.amazonaws.com

pulumi stack output rdsEndpoint
# Deve retornar: ecom-postgres.c1234567890ab.us-east-1.rds.amazonaws.com
```

### 3.4 Testar Aplicação
```bash
# Pegar DNS do Load Balancer
export LB_DNS=$(pulumi stack output loadBalancerDns)

# Testar se está respondendo
curl http://$LB_DNS/health
# Ou abra no navegador: http://$LB_DNS
```

---

## 🔐 Passo 4: Configurar GitHub Secrets

### 4.1 Acessar GitHub
1. Vá para seu repositório no GitHub
2. **Settings** → **Secrets and variables** → **Actions**

### 4.2 Criar Secrets
Clique em **New repository secret** para cada um:

#### Secret 1: AWS_ACCESS_KEY_ID
- Name: `AWS_ACCESS_KEY_ID`
- Value: Cole a access key do IAM user

#### Secret 2: AWS_SECRET_ACCESS_KEY
- Name: `AWS_SECRET_ACCESS_KEY`
- Value: Cole a secret access key do IAM user

#### Secret 3: PULUMI_ACCESS_TOKEN
- Name: `PULUMI_ACCESS_TOKEN`
- Value: Cole o token Pulumi gerado

#### Secret 4: PULUMI_STACK (Opcional)
- Name: `PULUMI_STACK`
- Value: `dev` (ou seu stack name)

### 4.3 Verificar Secrets
Na aba **Actions secrets and variables**, deve listar:
- ✅ AWS_ACCESS_KEY_ID
- ✅ AWS_SECRET_ACCESS_KEY
- ✅ PULUMI_ACCESS_TOKEN
- ✅ PULUMI_STACK

---

## 🔄 Passo 5: Deploy Automático via GitHub Actions

### 5.1 Fazer Commit das Mudanças
```bash
git add .github/workflows/ci-cd.yml
git commit -m "feat: enable AWS deployment in CI/CD pipeline"
git push origin main
```

### 5.2 Monitorar Pipeline
1. Acesse seu repositório no GitHub
2. Clique em **Actions**
3. Veja o workflow rodando
4. Acompanhe os 3 stages:
   - ✅ **Test & Lint** (2-3 min)
   - ✅ **Build Docker Image** (3-5 min)
   - ✅ **Deploy to AWS** (5-10 min)

### 5.3 Ver Logs
Clique no workflow em progresso → Clique em cada job para ver logs

### 5.4 Verificar Deploy
Após sucesso:
```bash
# Ver stacks no Pulumi
pulumi stack select dev

# Ver outputs
pulumi stack output

# Testar API
curl http://$(pulumi stack output loadBalancerDns)/health
```

---

## ⚙️ Passo 6: Configurar Aplicação em AWS

### 6.1 Variáveis de Ambiente no ECS

O deployment cria um ECS Fargate task. Precisa configurar variáveis de ambiente:

1. Acesse: https://console.aws.amazon.com/ecs/v2/
2. Procure: `service` (criado pelo Pulumi)
3. Clique em **Update service**
4. Vá em **Container details** → Clique no container `service`
5. Em **Environment**, adicione:

```
GOOSE_DBSTRING=host=RDS_ENDPOINT port=5432 user=postgres password=SENHA dbname=ecom sslmode=require
```

Onde:
- `RDS_ENDPOINT`: De `pulumi stack output rdsEndpoint`
- `SENHA`: A senha do RDS (configurada em `resources/rds.go`)

### 6.2 RDS Password

Atualmente no Pulumi está hardcoded. Para produção:

**NO ARQUIVO** `infra/pulumi/resources/rds.go`, mude de:
```go
MasterUserPassword: pulumi.String("password123"),
```

Para:
```go
MasterUserPassword: pulumi.StringPtr(
    config.GetSecret("rdsPassword"),
),
```

Depois configure:
```bash
cd infra/pulumi
pulumi config set --secret rdsPassword seu-super-seguro-password-aqui
pulumi up
```

---

## 🔍 Troubleshooting

### ❌ Erro: "ResourceNotFoundException" no RDS
**Causa**: RDS ainda está iniciando
**Solução**: Aguarde 5-10 minutos

### ❌ Erro: "Task failed to start"
**Causa**: Imagem Docker não existe ou variáveis de ambiente incorretas
**Solução**: 
```bash
# Verifique logs ECS
aws ecs describe-tasks --cluster cluster --tasks TASK_ARN
```

### ❌ Erro: "Access Denied" ao fazer push no ECR
**Causa**: Credenciais AWS incorretas
**Solução**: Verifique secrets no GitHub

### ❌ Pulumi "stack select" falha
**Causa**: Stack não existe
**Solução**:
```bash
cd infra/pulumi
pulumi stack ls  # Ver stacks
pulumi stack init dev  # Criar se não existe
```

### ❌ "Address already in use" na porta 5432 (tests locais)
**Causa**: PostgreSQL já rodando
**Solução**:
```bash
# Parar containers Docker
docker stop ecom-postgres
# Ou ver qual está usando:
lsof -i :5432
```

---

## 💰 Estimativa de Custos

| Recurso | Tipo | Custo/Mês (Aproximado) |
|---------|------|----------------------|
| VPC/Subnets | Free | $0 |
| Load Balancer | Network | $16.20 |
| ECS Fargate | 256 CPU + 512 MB | $5.63 |
| RDS PostgreSQL | db.t3.micro | $19.30 |
| NAT Gateway | Data Transfer | $32.00 |
| ECR Storage | Per GB | ~$5.00 |
| **TOTAL** | | **~$78-100/mês** |

**Nota**: Preços podem variar por região. Consulte AWS Pricing Calculator.

---

## 🛑 Destruir Recursos (Quando Quiser Parar)

### ⚠️ AVISO: Isso DELETA tudo!

```bash
cd infra/pulumi

# Ver o que será deletado
pulumi preview --destroy

# Deletar (pede confirmação)
pulumi destroy

# Deletar stack
pulumi stack rm dev
```

---

## 📚 Próximas Melhorias

- [ ] Adicionar Auto Scaling para ECS
- [ ] Implementar Blue/Green Deployment
- [ ] Adicionar CloudWatch Alarms
- [ ] Configurar RDS Backup automático
- [ ] Adicionar WAF para proteger Load Balancer
- [ ] Implementar Secrets Manager para senhas
- [ ] Adicionar GitOps com ArgoCD
- [ ] Implementar canary deployments

---

## 🔗 Links Úteis

- [AWS Console](https://console.aws.amazon.com/)
- [Pulumi App](https://app.pulumi.com/)
- [AWS ECR](https://console.aws.amazon.com/ecr/repositories)
- [AWS ECS](https://console.aws.amazon.com/ecs/v2/)
- [AWS RDS](https://console.aws.amazon.com/rds/)
- [Pulumi AWS Docs](https://www.pulumi.com/registry/packages/aws/)
- [AWS Pricing](https://calculator.aws/#/)

---

## ❓ FAQ

**P: Quanto tempo leva para fazer deploy?**
R: ~15-20 minutos total (testes + build + infra)

**P: Posso fazer rollback?**
R: Sim, via Pulumi: `pulumi refresh` → reverta código → `pulumi up`

**P: Posso usar diferente região?**
R: Sim, altere `aws:region` em Pulumi config

**P: Funciona com múltiplos ambientes?**
R: Sim! Crie stacks: `pulumi stack init prod`

**P: Como fazer manual deploy sem GitHub Actions?**
R: Execute `pulumi up` localmente após fazer commit

