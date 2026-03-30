# 📚 Documentação de Deploy - Índice Completo

## 🎯 Comece Aqui

### Para Deploy Rápido
👉 **[QUICK_START.md](./QUICK_START.md)** - 5 minutos para ter tudo rodando

### Para Entender Tudo
👉 **[AWS_DEPLOYMENT_GUIDE.md](./AWS_DEPLOYMENT_GUIDE.md)** - Guia completo passo-a-passo

### Para Executar Setup
👉 **[deploy-setup.sh](./deploy-setup.sh)** (Linux/Mac)  
👉 **[deploy-setup.bat](./deploy-setup.bat)** (Windows)

---

## 📁 Arquivos de Configuração

| Arquivo | Descrição |
|---------|-----------|
| `.github/workflows/ci-cd.yml` | Pipeline de CI/CD GitHub Actions |
| `IAM_POLICY.json` | Policy mínima de segurança para IAM user |
| `CI_CD_ANALYSIS.md` | Análise de possíveis pontos de erro |
| `DEPLOYMENT_CHECKLIST.md` | Checklist completo de deployment |

---

## 🏗️ Arquitetura na AWS

```
┌─────────────────────────────────────────┐
│         GitHub Actions CI/CD            │
│  (Tests → Build Docker → Deploy AWS)    │
└─────────────────┬───────────────────────┘
                  │
                  ▼
        ┌──────────────────────┐
        │  Amazon ECR          │
        │  (Container Registry)│
        └──────────────────────┘
                  │
                  ▼
    ┌─────────────────────────────────┐
    │    AWS VPC (us-east-1)          │
    │                                 │
    │  ┌─────────────────────────┐    │
    │  │  Public Subnets         │    │
    │  │  ┌─────────────────┐    │    │
    │  │  │ Network LB      │    │    │
    │  │  └────────┬────────┘    │    │
    │  └───────────┼──────────────┘    │
    │              │                   │
    │  ┌───────────▼──────────────┐    │
    │  │ Private Subnets          │    │
    │  │  ┌──────────────────┐    │    │
    │  │  │ ECS Fargate      │    │    │
    │  │  │ Container        │    │    │
    │  │  └────────┬─────────┘    │    │
    │  │           │              │    │
    │  │  ┌────────▼─────────┐    │    │
    │  │  │ RDS PostgreSQL   │    │    │
    │  │  └──────────────────┘    │    │
    │  └──────────────────────────┘    │
    └─────────────────────────────────┘
```

### Componentes Criados

1. **VPC** com subnets públicas e privadas
2. **Network Load Balancer** (porta 80/443)
3. **ECS Fargate** cluster para rodar containers
4. **ECR** repository para Docker images
5. **RDS PostgreSQL** database
6. **Security Groups** com regras apropriadas

---

## 📋 Checklist Rápido

- [ ] AWS CLI instalado
- [ ] Pulumi instalado
- [ ] Docker instalado
- [ ] AWS credentials configuradas
- [ ] Pulumi stack criado
- [ ] GitHub secrets configuradas
- [ ] Workflow testado

---

## 🔐 Segurança

### IAM User
- ✅ Use `IAM_POLICY.json` para limitar permissões
- ✅ Nunca compartilhe access keys
- ✅ Rotacione keys regularmente
- ✅ Use AWS Secrets Manager em produção

### Dados Sensíveis
- ✅ RDS password em variável de ambiente
- ✅ Secrets no GitHub Actions protegidas
- ✅ Pulumi token em variável de ambiente

---

## 💰 Custos

Estimativa mensal (região us-east-1):

| Recurso | Custo |
|---------|-------|
| Network Load Balancer | $16.20 |
| ECS Fargate (256 CPU, 512 MB) | $5.63 |
| RDS PostgreSQL (db.t3.micro) | $19.30 |
| NAT Gateway | $32.00 |
| ECR Storage | ~$5.00 |
| **TOTAL** | **~$78-100** |

*Para reduzir custos: desactive durante período sem uso*

---

## 🚀 Workflow Típico

### Desenvolvimento
```bash
git checkout -b feature/nova-funcionalidade
# Editar código
go test ./...
git commit -m "feat: nova funcionalidade"
git push origin feature/nova-funcionalidade
# Abrir Pull Request
```

### Merge para Main
```bash
# PR aprovado, merge para main
# GitHub Actions automáticamente:
# 1. Roda testes
# 2. Build Docker
# 3. Push para ECR
# 4. Deploy com Pulumi
```

### Monitorar Deploy
```bash
# Em GitHub Actions
github.com/SEU_REPO/actions

# Ou localmente
pulumi logs
pulumi stack output
```

---

## 🔄 Atualizações Comuns

### Atualizar Código
```bash
git add .
git commit -m "fix: bug tal"
git push origin main
# Deploy automático!
```

### Atualizar Infraestrutura
```bash
cd infra/pulumi
# Editar resources/
pulumi preview
pulumi up
# Ou push para main para CI/CD fazer
```

### Escalar Aplicação
```bash
# Editar infra/pulumi/resources/ecs.go
# Mudar DesiredCount: pulumi.Int(1) para pulumi.Int(3)
pulumi up
```

---

## 🐛 Debug

### Ver Logs da Aplicação
```bash
pulumi logs
```

### Ver Logs do ECS
```bash
aws logs tail /ecs/service --follow
```

### Conectar ao RDS
```bash
aws rds describe-db-instances --query 'DBInstances[0].[DBInstanceIdentifier,Endpoint.Address]'
# Depois conecte com psql ou similar
```

### Ver Estado Stack
```bash
pulumi stack select dev
pulumi stack output
```

---

## 📞 Suporte

### Links Importantes
- [AWS Console](https://console.aws.amazon.com/)
- [Pulumi Dashboard](https://app.pulumi.com/)
- [GitHub Actions](https://github.com/SEU_REPO/actions)
- [AWS Documentation](https://docs.aws.amazon.com/)
- [Pulumi Docs](https://www.pulumi.com/docs/)

### Comunidades
- [Pulumi Slack](https://slack.pulumi.com/)
- [AWS Forums](https://forums.aws.amazon.com/)
- [Stack Overflow](https://stackoverflow.com/questions/tagged/pulumi)

---

## ✅ Próximas Melhorias

- [ ] Auto Scaling Group
- [ ] Blue/Green Deployments
- [ ] GitOps com ArgoCD
- [ ] Multi-region
- [ ] Service Mesh
- [ ] Observability (DataDog, NewRelic)
- [ ] API Gateway + WAF
- [ ] Disaster Recovery

---

## 📝 Contribuindo

Encontrou um problema? Tem uma sugestão?
1. Abra uma issue em: `https://github.com/SEU_REPO/issues`
2. Descreva o problema
3. Envie PR com correção

---

## 📄 Licença

Mesmo projeto - Verifique LICENSE

---

## 🎉 Sucesso!

Seu projeto está pronto para production-grade deployment!

**Próximo passo**: Leia [QUICK_START.md](./QUICK_START.md)

