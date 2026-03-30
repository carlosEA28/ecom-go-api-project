# 🚀 Quick Start - Deploy em 5 Minutos

> Para deploy RÁPIDO na AWS. Versão detalhada em `AWS_DEPLOYMENT_GUIDE.md`

## Pré-requisitos (5 min)

```bash
# 1. Instale ferramentas
pip install awscli
# Windows: choco install pulumi
# Mac: brew install pulumi
# Linux: curl -fsSL https://get.pulumi.com | sh

# 2. Configure AWS
aws configure
# Cole: Access Key ID
# Cole: Secret Access Key
# Region: us-east-1

# 3. Teste conexão
aws sts get-caller-identity

# 4. Faça login no Pulumi
pulumi login
# Cole token de https://app.pulumi.com/account/tokens
```

## Setup (5 min)

```bash
cd infra/pulumi

# Crie stack
pulumi stack init dev

# Configure region
pulumi config set aws:region us-east-1

# Verifique preview
pulumi preview
```

## Deploy (10 min - Depende da AWS)

```bash
# Deploy tudo
pulumi up

# Copie os outputs:
# - loadBalancerDns
# - rdsEndpoint
# - ecrRepositoryUrl
```

## Configure GitHub (2 min)

1. Vá para: `https://github.com/SEU_REPO/settings/secrets/actions`
2. Clique **New repository secret** 4 vezes:

```
Name: AWS_ACCESS_KEY_ID
Value: Cole aqui sua access key

Name: AWS_SECRET_ACCESS_KEY
Value: Cole aqui sua secret

Name: PULUMI_ACCESS_TOKEN
Value: Cole seu token Pulumi

Name: PULUMI_STACK
Value: dev
```

3. Clique **Add secret**

## Push para Deploy (1 min)

```bash
git add .
git commit -m "feat: AWS deployment"
git push origin main

# Acesse GitHub Actions para monitorar:
# https://github.com/SEU_REPO/actions
```

## Pronto! 🎉

Seu app está deployado em AWS!

```bash
# Obter URL
pulumi stack output loadBalancerDns

# Testar
curl http://SEU_LOAD_BALANCER_DNS
```

---

## Troubleshooting Rápido

```bash
# Ver logs
pulumi logs

# Ver outputs
pulumi stack output

# Destruir (se errado)
pulumi destroy

# Refresh state
pulumi refresh
```

---

**Próximos passos**: Leia `AWS_DEPLOYMENT_GUIDE.md` para configuração em produção

