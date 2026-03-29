#!/bin/bash

# Script para personalizar o repositório clonado
# Este script remove referências ao repositório original e prepara tudo para você

set -e

# Cores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}ℹ${NC} $1"; }
log_success() { echo -e "${GREEN}✓${NC} $1"; }
log_warning() { echo -e "${YELLOW}⚠${NC} $1"; }
log_error() { echo -e "${RED}✗${NC} $1"; }

echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║   Personalizador de Repositório - Ecom Go API Project ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""

# 1. Perguntar informações
log_info "Vamos personalizar seu repositório"
echo ""

read -p "Seu GitHub username: " GITHUB_USERNAME
if [ -z "$GITHUB_USERNAME" ]; then
    log_error "Username é obrigatório"
    exit 1
fi

read -p "Nome do repositório [ecom-go-api-project]: " REPO_NAME
REPO_NAME=${REPO_NAME:-ecom-go-api-project}

read -p "Nome do projeto/módulo Go [ecom]: " GO_MODULE_NAME
GO_MODULE_NAME=${GO_MODULE_NAME:-ecom}

read -p "URL do novo repositório (deixe vazio para pular): " REPO_URL

echo ""
log_info "Configuração resumida:"
echo "  GitHub Username: $GITHUB_USERNAME"
echo "  Repositório: $REPO_NAME"
echo "  Módulo Go: $GO_MODULE_NAME"
echo "  URL (opcional): $REPO_URL"
echo ""

read -p "Continuar? (s/n): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Ss]$ ]]; then
    log_error "Abortado pelo usuário"
    exit 1
fi

echo ""
log_info "Iniciando personalização..."
echo ""

# 2. Remover histórico git
log_info "Removendo histórico git original..."
rm -rf .git
log_success "Histórico removido"

# 3. Atualizar go.mod
log_info "Atualizando go.mod..."
OLD_MODULE="github.com/sikozonpc/ecom"
NEW_MODULE="github.com/$GITHUB_USERNAME/$GO_MODULE_NAME"

sed -i.bak "s|$OLD_MODULE|$NEW_MODULE|g" go.mod
rm -f go.mod.bak

# Atualizar todos os imports no código
log_info "Atualizando imports em arquivos Go..."
find . -type f -name "*.go" -exec sed -i.bak "s|$OLD_MODULE|$NEW_MODULE|g" {} \;
find . -type f -name "*.go.bak" -delete

log_success "Módulo Go atualizado: $NEW_MODULE"

# 4. Atualizar main.go do Pulumi
log_info "Atualizando import no main.go do Pulumi..."
OLD_IMPORT="ecom-api/infra/pulumi/resources"
NEW_IMPORT="$NEW_MODULE/infra/pulumi/resources"

sed -i.bak "s|\"$OLD_IMPORT\"|\"$NEW_IMPORT\"|g" infra/pulumi/main.go
rm -f infra/pulumi/main.go.bak

log_success "Main.go do Pulumi atualizado"

# 5. Atualizar README
log_info "Atualizando README.md..."
if [ -f README.md ]; then
    sed -i.bak "s|sikozonpc/ecom|$GITHUB_USERNAME/$GO_MODULE_NAME|g" README.md
    sed -i.bak "s|sikozonpc|$GITHUB_USERNAME|g" README.md
    rm -f README.md.bak
    log_success "README.md atualizado"
else
    log_warning "README.md não encontrado"
fi

# 6. Atualizar docker-compose (se houver referências)
if [ -f docker-compose.yaml ]; then
    log_info "Verificando docker-compose.yaml..."
    if grep -q "sikozonpc" docker-compose.yaml; then
        sed -i.bak "s|sikozonpc|$GITHUB_USERNAME|g" docker-compose.yaml
        rm -f docker-compose.yaml.bak
        log_success "docker-compose.yaml atualizado"
    fi
fi

# 7. Criar novo repositório git local
log_info "Inicializando novo repositório git..."
git init
git config user.name "$GITHUB_USERNAME"
git config user.email "your-email@example.com"
git add .
git commit -m "Initial commit: Personalized from original ecom-go-api-project"
log_success "Repositório git criado"

# 8. Configurar remote (se URL foi fornecida)
if [ -n "$REPO_URL" ]; then
    log_info "Configurando remote do repositório..."
    git remote add origin "$REPO_URL"
    log_success "Remote 'origin' adicionado"
    echo ""
    log_info "Próximo passo - para fazer push:"
    echo "  git push -u origin main"
else
    log_warning "URL do repositório não fornecida"
    echo ""
    log_info "Próximos passos:"
    echo "  1. Crie um novo repositório no GitHub (github.com/$GITHUB_USERNAME/$REPO_NAME)"
    echo "  2. Execute: git remote add origin https://github.com/$GITHUB_USERNAME/$REPO_NAME.git"
    echo "  3. Execute: git push -u origin main"
fi

echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║              ✓ Personalização Completa!               ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""

log_success "Modificações aplicadas:"
echo "  ✓ Histórico git limpado"
echo "  ✓ Module path atualizado: $NEW_MODULE"
echo "  ✓ Imports Go atualizados"
echo "  ✓ Referências removidas"
echo "  ✓ Novo repositório git inicializado"
echo ""

log_info "Próximos passos:"
echo "  1. Ajuste credenciais de email no git:"
echo "     git config user.email 'seu-email@example.com'"
echo ""
echo "  2. Se ainda não fez push, crie o repo no GitHub e execute:"
echo "     git remote add origin https://github.com/$GITHUB_USERNAME/$REPO_NAME.git"
echo "     git branch -M main"
echo "     git push -u origin main"
echo ""
echo "  3. Configure os GitHub Secrets (Settings > Secrets):"
echo "     - AWS_ACCESS_KEY_ID"
echo "     - AWS_SECRET_ACCESS_KEY"
echo "     - RDS_PASSWORD"
echo ""

log_success "Pronto para começar! 🚀"
