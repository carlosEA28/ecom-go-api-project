# Database Migrations Guide

Este projeto usa **Goose** para gerenciar migrations do banco de dados PostgreSQL.

## 📋 Estrutura de Migrations

```
internal/adapters/postgresql/migrations/
├── 00001_create_products.sql
├── 00002_create_orders.sql
└── ...
```

Cada arquivo segue o padrão Goose com as seções:
- `-- +goose Up`: SQL para aplicar a migration
- `-- +goose Down`: SQL para reverter a migration

## 🚀 Como Usar

### 1. Instalando Goose Localmente

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 2. Variáveis de Ambiente Necessárias

```bash
export GOOSE_DBSTRING="host=localhost port=5432 user=postgres password=postgres dbname=ecom sslmode=disable"
```

### 3. Comandos Básicos

#### Aplicar todas as migrations pendentes
```bash
cd internal/adapters/postgresql
goose postgres "$GOOSE_DBSTRING" up
```

#### Ver status das migrations
```bash
cd internal/adapters/postgresql
goose postgres "$GOOSE_DBSTRING" status
```

#### Reverter a última migration
```bash
cd internal/adapters/postgresql
goose postgres "$GOOSE_DBSTRING" down
```

#### Reverter e reaplicar a última migration
```bash
cd internal/adapters/postgresql
goose postgres "$GOOSE_DBSTRING" redo
```

#### Executar até uma versão específica
```bash
cd internal/adapters/postgresql
goose postgres "$GOOSE_DBSTRING" up-to 1
```

## 📝 Criando Novas Migrations

### 1. Criar um novo arquivo SQL

```bash
cd internal/adapters/postgresql/migrations
touch 00003_create_orders_items.sql
```

### 2. Estrutura do arquivo

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_items (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id BIGINT NOT NULL REFERENCES products(id),
  quantity INTEGER NOT NULL CHECK (quantity > 0),
  price_at_purchase INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_items;
-- +goose StatementEnd
```

## 🔄 Fluxo no GitHub Actions

### CI/CD Pipeline Automático

1. **Push para main** → Pulumi cria RDS
2. **RDS fica pronto** → Goose aplica migrations automaticamente
3. **Migrations prontas** → ECS é deployado com código atualizado

### Migrations Manuais (GitHub Actions)

Use a workflow `migrations.yml` disparada manualmente:

```bash
# No GitHub, vá para: Actions > Database Migrations > Run workflow
# Selecione:
#   - Environment: dev/staging/prod
#   - Action: up/down/status/redo
```

## 🐳 Com Docker Compose

Se estiver desenvolvendo localmente com Docker:

```bash
# Subir PostgreSQL
docker-compose up -d postgres

# Aguarde o PostgreSQL ficar pronto (30 segundos)
sleep 30

# Aplicar migrations
cd internal/adapters/postgresql
goose postgres "host=localhost port=5432 user=postgres password=postgres dbname=ecom sslmode=disable" up

# Ver status
goose postgres "host=localhost port=5432 user=postgres password=postgres dbname=ecom sslmode=disable" status
```

## 🔐 Variáveis de Ambiente para RDS (AWS)

No arquivo `.env` ou GitHub Secrets:

```env
# Para desenvolvimento local com RDS
GOOSE_DBSTRING=host=seu-rds-endpoint.rds.amazonaws.com port=5432 user=postgres password=SUA_SENHA dbname=ecom sslmode=require

# Secrets do GitHub Actions
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- RDS_PASSWORD
```

## 🚨 Troubleshooting

### Problema: "no tables are present in the public schema"
**Solução**: Goose ainda não foi executado. Execute `goose up` para criar as tabelas.

### Problema: "connection refused"
**Solução**: Verifique se PostgreSQL está rodando e se as credenciais estão corretas.

### Problema: "migration lock"
**Solução**: Alguém está executando migrations em outro lugar. Se tiver certeza, limpe o lock:
```bash
goose postgres "$GOOSE_DBSTRING" fix
```

### Problema: "role 'postgres' does not exist"
**Solução**: Verifique o `RDS_USERNAME` nas variáveis de ambiente ou no Pulumi stack.

## 📚 Documentação

- [Goose Oficial](https://github.com/pressly/goose)
- [PostgreSQL](https://www.postgresql.org/docs/)
- [Pulumi RDS](https://www.pulumi.com/docs/reference/pkg/aws/rds/)

## ✅ Checklist para Deploy

- [ ] Código foi testado localmente com migrations
- [ ] Migrations foram validadas em staging
- [ ] GitHub Secrets estão configurados corretamente
- [ ] RDS está acessível de onde Goose vai rodar
- [ ] Pull request foi aprovado
- [ ] Deploy para main está pronto
