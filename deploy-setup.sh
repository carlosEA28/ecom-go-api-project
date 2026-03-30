#!/bin/bash
# Deploy Setup Script - Facilita o setup inicial para deploy na AWS
# Uso: ./deploy-setup.sh

set -e

echo "================================================"
echo "  🚀 E-Commerce API - AWS Deployment Setup"
echo "================================================"
echo ""

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check 1: AWS CLI
echo "${BLUE}[1/5] Verificando AWS CLI...${NC}"
if ! command -v aws &> /dev/null; then
    echo -e "${RED}✗ AWS CLI não instalado${NC}"
    echo "   Instale com: pip install awscli"
    exit 1
fi
echo -e "${GREEN}✓ AWS CLI encontrado$(aws --version)${NC}"

# Check 2: Pulumi
echo ""
echo "${BLUE}[2/5] Verificando Pulumi...${NC}"
if ! command -v pulumi &> /dev/null; then
    echo -e "${RED}✗ Pulumi não instalado${NC}"
    echo "   Windows: choco install pulumi"
    echo "   Mac: brew install pulumi"
    echo "   Linux: curl -fsSL https://get.pulumi.com | sh"
    exit 1
fi
echo -e "${GREEN}✓ Pulumi encontrado: $(pulumi version)${NC}"

# Check 3: Docker
echo ""
echo "${BLUE}[3/5] Verificando Docker...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker não instalado${NC}"
    echo "   Instale em: https://www.docker.com/products/docker-desktop"
    exit 1
fi
echo -e "${GREEN}✓ Docker encontrado: $(docker --version)${NC}"

# Check 4: Go
echo ""
echo "${BLUE}[4/5] Verificando Go...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go não instalado${NC}"
    exit 1
fi
GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✓ Go encontrado: ${GO_VERSION}${NC}"

# Check 5: AWS Credentials
echo ""
echo "${BLUE}[5/5] Verificando AWS Credentials...${NC}"
if ! aws sts get-caller-identity &> /dev/null; then
    echo -e "${RED}✗ AWS Credentials não configurados${NC}"
    echo "   Execute: aws configure"
    exit 1
fi
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
echo -e "${GREEN}✓ Credenciais AWS configuradas${NC}"
echo "   AWS Account: $ACCOUNT_ID"

# Tudo certo!
echo ""
echo "================================================"
echo -e "${GREEN}✓ Todas as dependências estão OK!${NC}"
echo "================================================"
echo ""

# Próximos passos
echo "${BLUE}📋 Próximos passos:${NC}"
echo ""
echo "1. Preparar Pulumi:"
echo "   cd infra/pulumi"
echo "   pulumi login"
echo "   pulumi stack init dev  (se não existe)"
echo "   pulumi config set aws:region us-east-1"
echo ""
echo "2. Preview do deploy:"
echo "   pulumi preview"
echo ""
echo "3. Fazer deploy (vai criar recursos AWS):"
echo "   pulumi up"
echo ""
echo "4. Configurar GitHub Secrets:"
echo "   https://github.com/SEU_USUARIO/SEU_REPO/settings/secrets/actions"
echo ""
echo "   Adicionar:"
echo "   - AWS_ACCESS_KEY_ID"
echo "   - AWS_SECRET_ACCESS_KEY"
echo "   - PULUMI_ACCESS_TOKEN"
echo "   - PULUMI_STACK"
echo ""
echo "5. Fazer push para main:"
echo "   git add ."
echo "   git commit -m 'feat: setup AWS deployment'"
echo "   git push origin main"
echo ""
echo "6. Monitorar deploy:"
echo "   GitHub → Actions → CI/CD Pipeline"
echo ""
echo "${YELLOW}Para mais detalhes, veja AWS_DEPLOYMENT_GUIDE.md${NC}"
