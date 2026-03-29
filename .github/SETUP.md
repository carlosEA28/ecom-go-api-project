# GitHub Actions Setup Guide

## 🔐 Configurar Secrets no GitHub

### Passo 1: Acessar Settings do Repositório

1. Vá para seu repositório no GitHub
2. Clique em **Settings**
3. No menu lateral, clique em **Secrets and variables** → **Actions**

### Passo 2: Adicionar Secrets Necessários

Clique em **New repository secret** e adicione cada um:

#### AWS Credentials
```
Name: AWS_ACCESS_KEY_ID
Value: sua-chave-de-acesso-aws
```

```
Name: AWS_SECRET_ACCESS_KEY
Value: sua-chave-secreta-aws
```

#### RDS Password
```
Name: RDS_PASSWORD
Value: postgres123!
```
(Deve ser a mesma definida em `infra/pulumi/resources/rds.go` linha 38)

#### Pulumi Stack (Opcional, se não usar S3)
```
Name: PULUMI_ACCESS_TOKEN
Value: seu-token-pulumi
```

### Passo 3: Estrutura de Secrets Recomendada

```
GitHub Secrets:
├── AWS_ACCESS_KEY_ID
├── AWS_SECRET_ACCESS_KEY
├── RDS_PASSWORD
└── PULUMI_ACCESS_TOKEN (opcional)
```

## 📋 Variáveis de Ambiente no Workflow

As seguintes variáveis estão configuradas nos workflows:

### `ci-cd.yml`
- `GO_VERSION`: 1.25.3
- `PULUMI_VERSION`: v3
- `AWS_REGION`: us-east-1

### `migrations.yml`
- `AWS_REGION`: us-east-1

## 🔑 Como Gerar Credenciais AWS

### 1. Access Key ID e Secret Access Key

```bash
# No AWS Console:
# 1. Vá para IAM > Users
# 2. Selecione seu usuário
# 3. Clique em "Create access key"
# 4. Escolha "Application running outside AWS"
# 5. Copie Access Key ID e Secret Access Key
```

### 2. Policy IAM Mínima Necessária

Crie uma policy com essas permissões:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchGetImage",
        "ecr:GetDownloadUrlForLayer",
        "ecr:PutImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload",
        "ecr:CreateRepository",
        "ecr:DeleteRepository",
        "ecr:DescribeRepositories"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ecs:UpdateService",
        "ecs:DescribeServices",
        "ecs:DescribeClusters"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "rds:DescribeDBInstances",
        "rds:DescribeDBClusters"
      ],
      "Resource": "*"
    }
  ]
}
```

## 🚀 Testando os Workflows

### Test Manual (GitHub Actions)

1. Vá para **Actions** no seu repositório
2. Selecione o workflow desejado
3. Clique em **Run workflow**
4. Escolha a branch e execute

### Local Testing

Para testar localmente antes de commitar:

```bash
# 1. Instalar act (simula GitHub Actions localmente)
# macOS
brew install act

# Linux
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | bash

# 2. Criar arquivo .env.local com seus secrets
cat > .env.local << EOF
AWS_ACCESS_KEY_ID=sua-chave
AWS_SECRET_ACCESS_KEY=sua-chave-secreta
RDS_PASSWORD=postgres123!
EOF

# 3. Executar workflow localmente
act -j build
```

## 📝 Estrutura do Workflow CI/CD

```
┌─────────────────────────────────────────────────┐
│           Push para main/develop                 │
└─────────────────┬───────────────────────────────┘
                  │
        ┌─────────▼────────────┐
        │   TEST JOB           │
        │ - go test ./...      │
        │ - PostgreSQL Service │
        └─────────┬────────────┘
                  │ (sucesso)
        ┌─────────▼────────────┐
        │   BUILD JOB          │
        │ - go build           │
        │ - Upload artifact    │
        └─────────┬────────────┘
                  │ (sucesso)
        ┌─────────▼────────────┐
        │   DEPLOY JOB         │
        │ (apenas main branch) │
        │ 1. Pulumi up         │
        │ 2. Goose migrations  │
        │ 3. Docker push       │
        │ 4. ECS update        │
        └─────────┬────────────┘
                  │ (completo)
        ┌─────────▼────────────┐
        │   DEPLOYMENT DONE    │
        │ ✓ Infraestrutura up  │
        │ ✓ DB migrado         │
        │ ✓ App deployado      │
        └──────────────────────┘
```

## 🔍 Monitorando Execução

### Ver logs da workflow

1. Vá para **Actions**
2. Clique na workflow que está rodando
3. Clique no job para ver logs em tempo real
4. Cada step mostra saída detalhada

### Troubleshooting comum

**Erro: "AWS credentials not found"**
- Verifique se os secrets estão no repositório correto
- Secrets não funcionam em forks (apenas no repo principal)

**Erro: "Pulumi stack not found"**
- Verifique se existe arquivo `Pulumi.dev.yaml` em `infra/pulumi`
- Verifique se a stack foi inicializada com `pulumi stack init`

**Erro: "RDS connection timeout"**
- RDS pode levar alguns minutos para ficar pronto
- O workflow aguarda até 5 minutos (ajustável)

## 🎯 Próximos Passos

1. ✅ Adicione todos os secrets no GitHub
2. ✅ Faça push para testar o workflow
3. ✅ Monitore os logs na aba Actions
4. ✅ Valide que infraestrutura foi criada no AWS
5. ✅ Verifique que migrations foram aplicadas

## 📞 Documentação Útil

- [GitHub Actions Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Pulumi GitHub Actions](https://github.com/pulumi/actions)
- [AWS Credentials Action](https://github.com/aws-actions/configure-aws-credentials)
