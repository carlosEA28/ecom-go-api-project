# Environment Variables Guide

Este documento descreve todas as variáveis de ambiente necessárias para executar o projeto em diferentes ambientes.

## 📋 Visão Geral

```
Desenvolvimento Local:    .env.local (gitignored)
Staging/Produção:        GitHub Secrets ou AWS Secrets Manager
CI/CD (GitHub Actions):  GitHub Secrets
Docker:                  docker-compose.yaml ou .env
```

---

## 🔧 Variáveis por Ambiente

### 1. Desenvolvimento Local

Arquivo: `.env.local` (ou `.env`)

```bash
# Database - Conexão local com PostgreSQL via Docker
GOOSE_DBSTRING=host=localhost port=5432 user=postgres password=postgres dbname=ecom sslmode=disable

# Aplicação
APP_PORT=8080
APP_ENV=development
LOG_LEVEL=debug

# AWS (opcional para local, necessário para deploy)
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=seu-access-key
AWS_SECRET_ACCESS_KEY=sua-secret-key

# Pulumi (para deployment)
PULUMI_STACK=dev
PULUMI_ACCESS_TOKEN=seu-token
```

### 2. GitHub Actions (CI/CD)

Local: **GitHub Secrets** → Settings > Secrets and variables > Actions

Secrets obrigatórios:
```
AWS_ACCESS_KEY_ID              → Chave de acesso AWS
AWS_SECRET_ACCESS_KEY          → Chave secreta AWS
RDS_PASSWORD                   → Senha do RDS (postgres123!)
PULUMI_ACCESS_TOKEN            → Token Pulumi (opcional se usar S3)
```

Variáveis no workflow (já configuradas):
```yaml
GO_VERSION: '1.25.3'
PULUMI_VERSION: 'v3'
AWS_REGION: 'us-east-1'
```

### 3. RDS (After Deployment)

Depois que Pulumi criar o RDS, obtenha do Pulumi stack:

```bash
cd infra/pulumi

# Obter automaticamente
pulumi stack output rdsEndpoint
pulumi stack output rdsPort
pulumi stack output rdsDatabase
pulumi stack output rdsUsername

# Ou configurar no .env
GOOSE_DBSTRING=host=seu-rds-endpoint.rds.amazonaws.com port=5432 user=postgres password=postgres123! dbname=ecom sslmode=require
```

---

## 📝 Descrição Detalhada de Cada Variável

### Database

#### `GOOSE_DBSTRING` (OBRIGATÓRIA)
**Descrição:** Connection string para PostgreSQL  
**Formato:** `host=X port=X user=X password=X dbname=X sslmode=X`  
**Desenvolvimento:** `host=localhost port=5432 user=postgres password=postgres dbname=ecom sslmode=disable`  
**Produção (RDS):** `host=seu-rds.rds.amazonaws.com port=5432 user=postgres password=SECURA dbname=ecom sslmode=require`  
**Usado por:** Goose (migrations), Aplicação (conexão BD)

---

### Application

#### `APP_PORT`
**Descrição:** Porta que a aplicação roda  
**Padrão:** `8080`  
**Desenvolvimento:** `8080`  
**Produção:** `8080` (ECS expõe via Load Balancer)

#### `APP_ENV`
**Descrição:** Ambiente de execução  
**Valores:** `development` | `staging` | `production`  
**Desenvolvimento:** `development`  
**Produção:** `production`

#### `LOG_LEVEL`
**Descrição:** Nível de logging  
**Valores:** `debug` | `info` | `warn` | `error`  
**Desenvolvimento:** `debug`  
**Produção:** `info` ou `warn`

---

### AWS (Infraestrutura e ECR)

#### `AWS_REGION`
**Descrição:** Região AWS  
**Padrão:** `us-east-1`  
**Valores comuns:** `us-east-1`, `us-west-2`, `eu-west-1`

#### `AWS_ACCESS_KEY_ID`
**Descrição:** ID da chave de acesso AWS  
**Onde obter:** AWS IAM → Users → Security credentials  
**Permissões necessárias:**
- `ecr:*` (ECR Repository)
- `ecs:UpdateService` (ECS Deploy)
- `rds:DescribeDBInstances` (RDS Info)
- `s3:*` (Pulumi Backend)

#### `AWS_SECRET_ACCESS_KEY`
**Descrição:** Chave secreta AWS  
**Onde obter:** AWS IAM → Users → Security credentials  
**⚠️ NUNCA compartilhe ou commite**

---

### RDS (Banco de Dados)

#### `RDS_ENDPOINT`
**Descrição:** Host do RDS  
**Formato:** `seu-database.xxxxx.rds.amazonaws.com`  
**Obtém de:** `pulumi stack output rdsEndpoint`

#### `RDS_PORT`
**Descrição:** Porta do RDS  
**Padrão:** `5432`  
**Obtém de:** `pulumi stack output rdsPort`

#### `RDS_DATABASE`
**Descrição:** Nome do database  
**Padrão:** `ecom`  
**Obtém de:** `pulumi stack output rdsDatabase`

#### `RDS_USERNAME`
**Descrição:** Usuário do RDS  
**Padrão:** `postgres`  
**Obtém de:** `pulumi stack output rdsUsername`

#### `RDS_PASSWORD`
**Descrição:** Senha do RDS  
**Valor:** Definida em `infra/pulumi/resources/rds.go:136`  
**Padrão:** `postgres123!`  
**⚠️ NUNCA commite, use GitHub Secrets**

---

### Pulumi (Infrastructure as Code)

#### `PULUMI_STACK`
**Descrição:** Stack do Pulumi a usar  
**Valores:** `dev` | `staging` | `prod`  
**Padrão:** `dev`

#### `PULUMI_ACCESS_TOKEN`
**Descrição:** Token para acessar Pulumi backend  
**Onde obter:** https://app.pulumi.com → Account settings → Access tokens  
**Alternativa:** Usar S3 backend com `PULUMI_CLOUD_URL`

#### `PULUMI_CLOUD_URL`
**Descrição:** URL do backend Pulumi (S3)  
**Formato:** `s3://seu-bucket-name`  
**Quando usar:** Se não tiver account Pulumi

---

## 🚀 Setup por Cenário

### Cenário 1: Desenvolvimento Local

```bash
# 1. Criar .env.local (gitignored)
cat > .env.local << EOF
GOOSE_DBSTRING=host=localhost port=5432 user=postgres password=postgres dbname=ecom sslmode=disable
APP_PORT=8080
APP_ENV=development
LOG_LEVEL=debug
EOF

# 2. Exportar variáveis
export $(cat .env.local | xargs)

# 3. Iniciar PostgreSQL
docker-compose up -d postgres

# 4. Aplicar migrations
./migrate.sh dev up

# 5. Rodar aplicação
go run ./cmd
```

### Cenário 2: Preparar Deploy (GitHub Actions)

```bash
# 1. Adicionar secrets no GitHub
gh secret set AWS_ACCESS_KEY_ID
gh secret set AWS_SECRET_ACCESS_KEY
gh secret set RDS_PASSWORD

# 2. Configurar Pulumi stack
cd infra/pulumi
pulumi stack init dev
pulumi config set aws:region us-east-1

# 3. Fazer push
git add .
git commit -m "Setup for deployment"
git push origin main

# 4. GitHub Actions executa automaticamente
# - Tests
# - Build
# - Deploy (Pulumi up)
# - Migrations (Goose up)
```

### Cenário 3: Usar RDS Existente

```bash
# 1. Obter informações do RDS
cd infra/pulumi
RDS_ENDPOINT=$(pulumi stack output rdsEndpoint)
RDS_PASSWORD=$(aws rds describe-db-instances --query 'DBInstances[0].MasterUserPassword')

# 2. Criar .env.local com essas infos
cat > .env.local << EOF
GOOSE_DBSTRING=host=$RDS_ENDPOINT port=5432 user=postgres password=$RDS_PASSWORD dbname=ecom sslmode=require
APP_ENV=staging
LOG_LEVEL=info
EOF

# 3. Aplicar migrations
./migrate.sh staging up

# 4. Conectar aplicação
go run ./cmd
```

---

## 📊 Matriz de Variáveis por Ambiente

| Variável | Local | Dev | Staging | Prod | GitHub Actions |
|----------|-------|-----|---------|------|----------------|
| GOOSE_DBSTRING | ✅ local | ✅ RDS | ✅ RDS | ✅ RDS | - |
| APP_PORT | 8080 | 8080 | 8080 | 8080 | - |
| APP_ENV | development | development | staging | production | - |
| LOG_LEVEL | debug | info | info | warn | - |
| AWS_REGION | ⚠️ optional | ✅ | ✅ | ✅ | ✅ (Secrets) |
| AWS_ACCESS_KEY_ID | ⚠️ optional | ✅ | ✅ | ✅ | ✅ (Secrets) |
| AWS_SECRET_ACCESS_KEY | ⚠️ optional | ✅ | ✅ | ✅ | ✅ (Secrets) |
| RDS_ENDPOINT | - | ✅ pulumi output | ✅ | ✅ | - |
| RDS_PASSWORD | - | - | - | ✅ | ✅ (Secrets) |
| PULUMI_STACK | dev | dev | staging | prod | - |
| PULUMI_ACCESS_TOKEN | ⚠️ optional | ⚠️ optional | ⚠️ optional | ⚠️ optional | ✅ (Secrets) |

---

## ⚙️ Como Usar Em Diferentes Cenários

### Docker Compose

```bash
# .env ou .env.local
GOOSE_DBSTRING=host=postgres port=5432 user=postgres password=postgres dbname=ecom sslmode=disable

# docker-compose.yaml já lê automaticamente
docker-compose up
```

### GitHub Actions

```yaml
# Variáveis de ambiente no workflow
env:
  AWS_REGION: us-east-1

# Secrets são automaticamente disponíveis
# - AWS_ACCESS_KEY_ID
# - AWS_SECRET_ACCESS_KEY
# - RDS_PASSWORD
```

### AWS Lambda (Future)

Se mudar para Lambda, use AWS Secrets Manager:

```bash
aws secretsmanager create-secret \
  --name ecom/production \
  --secret-string '{"db_host":"...","db_password":"..."}'
```

---

## 🔐 Security Best Practices

### ✅ DO

- ✅ Use `.env.local` para desenvolvimento
- ✅ Use GitHub Secrets para CI/CD
- ✅ Use AWS Secrets Manager para produção
- ✅ Rotacione passwords regularmente
- ✅ Use `.gitignore` para `.env`
- ✅ Use `sslmode=require` em produção

### ❌ DON'T

- ❌ NÃO commite `.env` com valores reais
- ❌ NÃO compartilhe `AWS_SECRET_ACCESS_KEY`
- ❌ NÃO use `sslmode=disable` em produção
- ❌ NÃO use a mesma senha em todos ambientes
- ❌ NÃO commite `PULUMI_ACCESS_TOKEN`

---

## 🆘 Troubleshooting

### Erro: "connection refused"
```bash
# Problema: PostgreSQL não está rodando
# Solução:
docker-compose up -d postgres
sleep 30
./migrate.sh dev up
```

### Erro: "permission denied" no RDS
```bash
# Problema: Credenciais incorretas
# Solução:
pulumi stack output rdsUsername
pulumi stack output rdsEndpoint
# Verifique a senha no arquivo rds.go
```

### Erro: "secret not found" (GitHub Actions)
```bash
# Problema: Secrets não configurados
# Solução:
gh secret set AWS_ACCESS_KEY_ID
gh secret set AWS_SECRET_ACCESS_KEY
gh secret set RDS_PASSWORD
```

---

## 📚 Referência Rápida

```bash
# Obter valores do Pulumi
cd infra/pulumi
pulumi stack output rdsEndpoint
pulumi stack output rdsPort
pulumi stack output rdsDatabase
pulumi stack output rdsUsername

# Adicionar secrets
gh secret set <NOME>
gh secret list

# Testar conexão
psql "host=localhost port=5432 user=postgres password=postgres dbname=ecom sslmode=disable" -c "SELECT 1"

# Ver variáveis ativas
env | grep -E "GOOSE|APP_|AWS_|RDS_|PULUMI_"
```

---

## ✅ Checklist de Setup

- [ ] Copiar `.env.example` para `.env.local`
- [ ] Editar `.env.local` com seus valores
- [ ] Adicionar `.env.local` ao `.gitignore`
- [ ] Testar conexão: `psql $GOOSE_DBSTRING -c "SELECT 1"`
- [ ] Configurar GitHub Secrets: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `RDS_PASSWORD`
- [ ] Testar localmente: `go run ./cmd`
- [ ] Testar migrations: `./migrate.sh dev up`
- [ ] Fazer push e monitorar GitHub Actions

Pronto! 🚀
