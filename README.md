# 🛒 E-Commerce API - Go

API de e-commerce desenvolvida em Go com PostgreSQL, Docker e deploy automático na AWS.

## 📁 Estrutura do Projeto

```
ecom-go-api-project/
├── cmd/                          # Aplicação principal
│   ├── main.go                  # Entry point
│   └── api.go                   # Setup da API
├── internal/                     # Código interno
│   ├── adapters/                # Adaptadores (PostgreSQL)
│   ├── products/                # Domain de produtos
│   ├── orders/                  # Domain de pedidos
│   ├── env/                     # Variáveis de ambiente
│   └── json/                    # Helpers JSON
├── infra/                        # Infraestrutura
│   └── pulumi/                  # IaC com Pulumi (AWS)
├── .github/                      # GitHub Actions
│   └── workflows/
│       └── ci-cd.yml           # Pipeline de CI/CD
├── deployment-docs/             # 📚 Documentação de deploy
│   ├── QUICK_START.md          # Start rápido (5 min)
│   ├── AWS_DEPLOYMENT_GUIDE.md # Guia detalhado
│   ├── ACTION_MAP.md           # Mapa visual
│   ├── deploy-setup.sh/.bat    # Scripts de setup
│   └── ...
├── Dockerfile                    # Build container
├── docker-compose.yaml          # PostgreSQL local
└── go.mod/go.sum               # Dependencies
```

## 🚀 Quick Start

### Desenvolvimento Local

```bash
# 1. Start PostgreSQL
docker-compose up -d

# 2. Run migrations
./migrate.sh  # Linux/Mac
# ou
migrate.bat   # Windows

# 3. Run aplicação
go run ./cmd

# 4. API disponível em
curl http://localhost:8080
```

### Deploy na AWS

Para fazer deploy na AWS com CI/CD automático:

**→ [Leia deployment-docs/QUICK_START.md](./deployment-docs/QUICK_START.md)**

(5 minutos para ter tudo rodando)

## 📚 Documentação

Toda a documentação de deploy está em `deployment-docs/`:

| Arquivo | Descrição |
|---------|-----------|
| **QUICK_START.md** | Setup em 5 minutos |
| **AWS_DEPLOYMENT_GUIDE.md** | Guia passo-a-passo completo |
| **ACTION_MAP.md** | Mapa visual de decisões |
| **CI_CD_ANALYSIS.md** | Análise de possíveis erros |
| **DEPLOYMENT_CHECKLIST.md** | Checklist de validação |
| **deploy-setup.sh/.bat** | Scripts de verificação |

## 🏗️ Arquitetura

### Local
- Go API na porta 8080
- PostgreSQL na porta 5432
- Docker Compose para ambiente

### AWS (via Pulumi + GitHub Actions)
- VPC com subnets públicas/privadas
- Network Load Balancer
- ECS Fargate para containers
- RDS PostgreSQL
- ECR para Docker images
- Security Groups

## 🔄 CI/CD Pipeline

Automático ao fazer push para `main`:

1. **Test & Lint** - Testes Go + linting
2. **Build Docker** - Build imagem + push ECR
3. **Deploy AWS** - Pulumi faz deploy automático

Ver: `.github/workflows/ci-cd.yml`

## 🛠️ Tecnologias

- **Go 1.25.6** - Linguagem
- **PostgreSQL 16** - Banco de dados
- **Docker** - Containerização
- **Pulumi** - Infrastructure as Code
- **GitHub Actions** - CI/CD
- **AWS** - Cloud provider

## 📝 Comandos Úteis

```bash
# Testes
go test ./...

# Build local
docker build -t ecom-api:latest .

# Linting
golangci-lint run

# Formato
gofmt -s -w .

# Deploy Pulumi
cd infra/pulumi
pulumi preview
pulumi up

# Ver logs Pulumi
pulumi logs
```

## 🔐 Configuração

### Variáveis de Ambiente

```bash
# Database
GOOSE_DBSTRING=host=localhost port=5432 user=postgres password=postgres dbname=ecom sslmode=disable

# AWS (Production)
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=xxx
AWS_SECRET_ACCESS_KEY=xxx
```

Ver `.env.example` para mais detalhes.

## ⚠️ Notas Importantes

- Senhas em `.env` não devem ser commitadas
- Use AWS Secrets Manager em produção
- Verifique `deployment-docs/CI_CD_ANALYSIS.md` para problemas conhecidos

## 📞 Suporte

- Documentação de deploy: `deployment-docs/`
- AWS: https://docs.aws.amazon.com/
- Pulumi: https://www.pulumi.com/docs/
- Go: https://golang.org/doc/

## 📄 Licença

[Especificar sua licença aqui]

---

**Pronto para fazer deploy?** → [deployment-docs/QUICK_START.md](./deployment-docs/QUICK_START.md) 🚀

