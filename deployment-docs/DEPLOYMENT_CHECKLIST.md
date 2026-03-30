# 📋 Checklist de Deploy - AWS

## ✅ Fase 1: Setup Local (2-3 horas)

### Ferramentas
- [ ] **AWS CLI** instalado (`aws --version`)
- [ ] **Pulumi** instalado (`pulumi version`)
- [ ] **Docker** instalado (`docker --version`)
- [ ] **Go 1.25.6+** instalado (`go version`)

### AWS Account
- [ ] Conta AWS criada
- [ ] IAM User criado (nome: `github-deployer`)
- [ ] Access Key gerada e guardada (NUNCA COMPARTILHE!)
- [ ] AWS CLI configurado (`aws configure`)
- [ ] Credenciais testadas (`aws sts get-caller-identity`)

### Pulumi Setup
- [ ] Conta Pulumi criada (https://app.pulumi.com/signup)
- [ ] Pulumi Access Token gerado
- [ ] `pulumi login` executado
- [ ] Stack `dev` criado (`pulumi stack init dev`)
- [ ] AWS Region configurada (`pulumi config set aws:region us-east-1`)

### Teste Local
- [ ] `pulumi preview` sem erros
- [ ] Recursos AWS no preview fazem sentido
- [ ] Custo estimado dentro do orçamento

---

## ✅ Fase 2: Deploy Piloto (1-2 horas)

### Deploy Inicial
- [ ] `pulumi up` executado com sucesso
- [ ] Todos os recursos criados:
  - [ ] VPC com Subnets (públicas e privadas)
  - [ ] Security Groups (ALB, ECS, RDS)
  - [ ] RDS PostgreSQL
  - [ ] ECR Repository
  - [ ] ECS Fargate Cluster
  - [ ] Network Load Balancer
- [ ] Outputs capturados (`pulumi stack output`)

### Validação
- [ ] Load Balancer DNS acessível
- [ ] RDS endpoint no output correto
- [ ] ECR repository criado e acessível
- [ ] Teste de conectividade ao RDS (opcional)

---

## ✅ Fase 3: GitHub Integration (30 minutos)

### Secrets no GitHub
1. Acesse: `https://github.com/SEU_USUARIO/SEU_REPO/settings/secrets/actions`
2. Crie secret `AWS_ACCESS_KEY_ID`
   - [ ] Valor copiado do IAM user
3. Crie secret `AWS_SECRET_ACCESS_KEY`
   - [ ] Valor copiado do IAM user
4. Crie secret `PULUMI_ACCESS_TOKEN`
   - [ ] Valor copiado de Pulumi
5. Crie secret `PULUMI_STACK` (opcional)
   - [ ] Valor: `dev`

### Verificação
- [ ] Todos os 4 secrets aparecem na lista
- [ ] Nenhum valor contém espaços extras

---

## ✅ Fase 4: CI/CD Pipeline (1 hora)

### Commit Workflow Atualizado
```bash
git add .github/workflows/ci-cd.yml
git commit -m "feat: activate AWS deployment in CI/CD"
git push origin main
```

- [ ] Commit feito
- [ ] Push para main concluído

### Monitorar Primeira Execução
1. Acesse: `https://github.com/SEU_USUARIO/SEU_REPO/actions`
2. Monitore os 3 stages:
   - [ ] **Test & Lint** (passar)
   - [ ] **Build Docker** (push para GHCR)
   - [ ] **Deploy to AWS** (Pulumi up)

### Validar Deploy
```bash
cd infra/pulumi
pulumi stack output loadBalancerDns
# Deve retornar: ecom-nlb-xxxxxx.elb.us-east-1.amazonaws.com

curl http://LOAD_BALANCER_DNS/health
# Deve retornar: 200 OK ou mensagem de health
```

- [ ] Aplicação respondendo no ALB
- [ ] Logs do ECS sem erros
- [ ] RDS com dados

---

## ✅ Fase 5: Produção (Configuração)

### Variáveis de Ambiente
```bash
# Pegar RDS Endpoint:
pulumi stack output rdsEndpoint

# Configurar no ECS (via AWS Console ou CLI):
export GOOSE_DBSTRING="host=RDS_ENDPOINT port=5432 user=postgres password=SENHA dbname=ecom sslmode=require"
```

- [ ] `GOOSE_DBSTRING` setada corretamente
- [ ] Migrations executadas
- [ ] Conexão ao RDS verificada

### Segurança
- [ ] Senhas RDS movidas para AWS Secrets Manager
- [ ] IAM policies restringidas ao mínimo necessário
- [ ] Bucket S3 para backups criado
- [ ] VPC endpoints configurados (opcional)
- [ ] WAF habilitado no ALB (opcional)

---

## ✅ Fase 6: Monitoramento

### CloudWatch
- [ ] ECS Cluster logs visíveis
- [ ] RDS Performance Insights habilitado
- [ ] ALB logs ativados

### Alertas
- [ ] High CPU > 80%
- [ ] High Memory > 80%
- [ ] Failed Tasks > 0
- [ ] RDS Storage > 80%

### Backups
- [ ] RDS automated backups: 7 dias
- [ ] ECR image cleanup policy criada
- [ ] Stack configs fazer backup

---

## 🎯 Rotina Diária

### Deploy
```bash
# Local - antes de push
./deploy-setup.sh  # ou .bat no Windows
pulumi preview
pulumi up --refresh

# Automático via GitHub
git add .
git commit -m "feat: new feature"
git push origin main
# Monitore em GitHub Actions
```

### Verificações
```bash
# Health check
curl http://LOAD_BALANCER_DNS/health

# Logs
pulumi logs

# Stack state
pulumi stack output
```

### Updates
```bash
# Atualizar dependências Go
go get -u

# Rebuild local
docker build -t myapp:test .

# Deploy
git push origin main
```

---

## ⚠️ Troubleshooting Rápido

| Problema | Solução |
|----------|---------|
| `Error: 403 Unauthorized` no ECR | Verifique AWS credentials, execute `aws configure` |
| `stack not found` | Execute `pulumi stack init dev` |
| `port 5432 already in use` | `docker stop ecom-postgres` |
| `GitHub Action falha` | Verifique secrets em Settings → Secrets |
| `RDS connection timeout` | Aguarde 5-10 min, RDS pode estar iniciando |
| `ECS task failed` | Verifique logs: `aws ecs describe-tasks --cluster ... --tasks ...` |

---

## 💾 Backups e Recovery

### Fazer Backup do Stack
```bash
# Exportar estado
pulumi stack export > stack-backup.json

# Também fazer backup local
git add .
git commit -m "backup: stack state"
```

### Recuperar de Falha
```bash
# Se houver erro, voltar ao estado anterior
git checkout HEAD~1
pulumi refresh  # Sincroniza estado local com AWS
pulumi up --refresh
```

### Destruir e Reconstruir
```bash
# Destruir tudo (CUIDADO!)
pulumi destroy

# Reconstruir
pulumi up
```

---

## 📊 Dashboard Pulumi

Monitore tudo em: https://app.pulumi.com/

- [ ] Stack `dev` visível
- [ ] Últimas operações listadas
- [ ] Nenhum erro em vermelho

---

## 🎓 Aprendizado Contínuo

### Documentação
- [ ] Ler: [Pulumi AWS Guide](https://www.pulumi.com/docs/clouds/aws/get-started/)
- [ ] Ler: [ECS Best Practices](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/Welcome.html)
- [ ] Ler: [RDS Best Practices](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Welcome.html)

### Melhorias Futuras
- [ ] [ ] Auto Scaling Group para ECS
- [ ] [ ] Blue/Green Deployments
- [ ] [ ] Canary Deployments
- [ ] [ ] Multi-region setup
- [ ] [ ] GitOps com ArgoCD
- [ ] [ ] Service Mesh (Istio/Linkerd)

---

## 📞 Suporte

### Documentação
- AWS: https://docs.aws.amazon.com/
- Pulumi: https://www.pulumi.com/docs/
- Go: https://pkg.go.dev/

### Comunidades
- Pulumi Slack: https://slack.pulumi.com/
- AWS Forums: https://forums.aws.amazon.com/
- Stack Overflow: Tag `pulumi`, `aws`, `go`

---

## 📝 Notas

**Data de Início**: ___________
**Data de Conclusão**: ___________
**Ambiente**: ☐ Dev ☐ Staging ☐ Production
**Pessoa Responsável**: ___________________

**Observações**:
```
________________________________
________________________________
________________________________
```

---

## ✨ Sucesso!

Quando tudo estiver funcionando:
1. Fazer commit celebrando: `git commit -m "🎉 deployment success!"`
2. Notificar time
3. Documentar lições aprendidas
4. Celebrar! 🎉

