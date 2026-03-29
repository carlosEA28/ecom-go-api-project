#!/bin/bash

# Script para executar migrations localmente ou em RDS
# Uso: ./migrate.sh [dev|staging|prod] [up|down|status|redo]

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para imprimir mensagens coloridas
log_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

log_success() {
    echo -e "${GREEN}✓${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Validar argumentos
if [ $# -lt 1 ]; then
    log_error "Uso: ./migrate.sh [dev|staging|prod] [up|down|status|redo]"
    echo ""
    echo "Exemplos:"
    echo "  ./migrate.sh dev                 # Apply migrations (padrão)"
    echo "  ./migrate.sh dev up              # Apply all pending migrations"
    echo "  ./migrate.sh dev down            # Revert last migration"
    echo "  ./migrate.sh dev status          # Check migration status"
    echo "  ./migrate.sh dev redo            # Revert and reapply last"
    exit 1
fi

ENVIRONMENT=${1:-dev}
ACTION=${2:-up}

# Validar ambiente
case "$ENVIRONMENT" in
    dev|staging|prod)
        log_info "Environment: $ENVIRONMENT"
        ;;
    *)
        log_error "Ambiente inválido: $ENVIRONMENT (use: dev, staging, prod)"
        exit 1
        ;;
esac

# Validar ação
case "$ACTION" in
    up|down|status|redo)
        log_info "Action: $ACTION"
        ;;
    *)
        log_error "Ação inválida: $ACTION (use: up, down, status, redo)"
        exit 1
        ;;
esac

# Verificar se goose está instalado
if ! command -v goose &> /dev/null; then
    log_warning "Goose não está instalado"
    log_info "Instalando Goose..."
    go install github.com/pressly/goose/v3/cmd/goose@latest
fi

# Determinar se é local ou RDS
if [ "$ENVIRONMENT" = "dev" ]; then
    log_info "Usando PostgreSQL local (docker-compose)"
    
    # Verificar se docker-compose está rodando
    if ! docker ps | grep -q postgres; then
        log_warning "PostgreSQL não está rodando no Docker"
        log_info "Iniciando docker-compose..."
        docker-compose up -d postgres
        sleep 10
    fi
    
    DB_HOST="localhost"
    DB_PORT="5432"
    DB_USER="postgres"
    DB_PASSWORD="postgres"
    DB_NAME="ecom"
    DB_SSLMODE="disable"
else
    log_info "Obtendo credenciais do Pulumi para $ENVIRONMENT..."
    
    # Verificar se pulumi está instalado
    if ! command -v pulumi &> /dev/null; then
        log_error "Pulumi não está instalado"
        exit 1
    fi
    
    cd infra/pulumi
    
    # Selecionar stack
    pulumi stack select "$ENVIRONMENT" 2>/dev/null || {
        log_error "Stack '$ENVIRONMENT' não existe"
        exit 1
    }
    
    # Obter valores do stack
    DB_HOST=$(pulumi stack output rdsEndpoint 2>/dev/null)
    DB_PORT=$(pulumi stack output rdsPort 2>/dev/null)
    DB_USER=$(pulumi stack output rdsUsername 2>/dev/null)
    DB_NAME=$(pulumi stack output rdsDatabase 2>/dev/null)
    DB_SSLMODE="require"
    
    if [ -z "$RDS_PASSWORD" ]; then
        read -sp "Enter RDS password: " DB_PASSWORD
        echo ""
    else
        DB_PASSWORD="$RDS_PASSWORD"
    fi
    
    cd - > /dev/null
fi

# Construir connection string
DB_STRING="host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=$DB_SSLMODE"

log_info "Conectando a: $DB_HOST:$DB_PORT/$DB_NAME"

# Testar conexão
if ! psql "$DB_STRING" -c "SELECT 1" &> /dev/null; then
    log_error "Não conseguiu conectar ao banco de dados"
    exit 1
fi

log_success "Conectado ao banco de dados"

# Executar migrations
cd internal/adapters/postgresql
log_info "Executando migrations..."

case "$ACTION" in
    up)
        log_info "⬆️  Aplicando todas as migrations pendentes..."
        if goose postgres "$DB_STRING" up; then
            log_success "Migrations aplicadas com sucesso"
        else
            log_error "Erro ao aplicar migrations"
            exit 1
        fi
        ;;
    down)
        log_info "⬇️  Revertendo última migration..."
        if goose postgres "$DB_STRING" down; then
            log_success "Migration revertida com sucesso"
        else
            log_error "Erro ao reverter migration"
            exit 1
        fi
        ;;
    status)
        log_info "📊 Status das migrations:"
        echo ""
        goose postgres "$DB_STRING" status
        ;;
    redo)
        log_info "🔄 Revertendo e reaplicando última migration..."
        if goose postgres "$DB_STRING" redo; then
            log_success "Migration redo completa"
        else
            log_error "Erro ao fazer redo"
            exit 1
        fi
        ;;
esac

cd - > /dev/null

# Mostrar status final
log_success "Operação concluída!"
echo ""
log_info "Resumo:"
echo "  Ambiente: $ENVIRONMENT"
echo "  Host: $DB_HOST:$DB_PORT"
echo "  Database: $DB_NAME"
echo "  Ação: $ACTION"
echo ""

# Status final
cd internal/adapters/postgresql
log_info "Status atual das migrations:"
goose postgres "$DB_STRING" status
cd - > /dev/null
