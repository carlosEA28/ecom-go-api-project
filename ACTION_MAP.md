# 🗺️ Mapa de Ação - Deploy na AWS

## 👤 ESTOU AQUI?

### ▶️ SIM - Sou iniciante, quero começar já
**→ Leia:** [QUICK_START.md](./QUICK_START.md) (5 minutos)  
**→ Depois:** [AWS_DEPLOYMENT_GUIDE.md](./AWS_DEPLOYMENT_GUIDE.md) (30 minutos)

### ▶️ SIM - Quero entender toda a arquitetura
**→ Leia:** [DEPLOYMENT_DOCS_INDEX.md](./DEPLOYMENT_DOCS_INDEX.md)  
**→ Depois:** [CI_CD_ANALYSIS.md](./CI_CD_ANALYSIS.md)

### ▶️ SIM - Já tenho conta AWS, só quero fazer deploy
**→ Execute:** `./deploy-setup.sh` (Linux/Mac) ou `deploy-setup.bat` (Windows)  
**→ Depois:** Siga [AWS_DEPLOYMENT_GUIDE.md](./AWS_DEPLOYMENT_GUIDE.md) - Passo 2 em diante

### ▶️ SIM - Sou DevOps/SRE experiente
**→ Leia:** [CI_CD_ANALYSIS.md](./CI_CD_ANALYSIS.md) (problema já identificados)  
**→ Analise:** `.github/workflows/ci-cd.yml`  
**→ Customize:** Conforme sua infra existente

---

## 📊 Tempo Estimado por Tarefa

```
┌─────────────────────────────────────────────────┐
│ SETUP INICIAL                                   │
│ ├─ Instalar ferramentas ......... 10 min        │
│ ├─ Criar AWS Account ............ 15 min        │
│ ├─ Configurar Pulumi ............ 10 min        │
│ └─ SUBTOTAL ..................... 35 min        │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│ DEPLOY INICIAL                                  │
│ ├─ Pulumi preview ............... 2 min         │
│ ├─ Pulumi up (criar infraestrutura) 15 min     │
│ └─ SUBTOTAL ..................... 17 min        │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│ CI/CD SETUP                                     │
│ ├─ Configurar GitHub secrets ..... 5 min        │
│ ├─ Fazer commit .................. 2 min        │
│ ├─ GitHub Actions rodando ........ 15 min       │
│ └─ SUBTOTAL ..................... 22 min        │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│ TOTAL PARA PRODUÇÃO: ~75 minutos (1h15m)       │
└─────────────────────────────────────────────────┘
```

---

## 🔄 Fluxo de Decisão

```
                        COMECE AQUI
                             │
                             ▼
                   ┌─────────────────────┐
                   │ Tem AWS Account?    │
                   └─────────────────────┘
                      NÃO ↓  ↓ SIM
                          │  │
            ┌─────────────┘  └─────────────┐
            ▼                               ▼
    ┌──────────────┐              ┌──────────────┐
    │ Criar AWS    │              │ Tem AWS CLI? │
    │ Account      │              └──────────────┘
    │ (15 min)     │                 NÃO ↓  ↓ SIM
    └──────────────┘                    │  │
            │            ┌──────────────┘  │
            │            ▼                  │
            │    ┌──────────────┐          │
            │    │ Instalar     │          │
            │    │ AWS CLI      │          │
            │    │ (10 min)     │          │
            │    └──────────────┘          │
            │            │                  │
            └────────────┼──────────────────┘
                         ▼
                ┌──────────────────┐
                │ aws configure    │
                │ (5 min)          │
                └──────────────────┘
                         │
                         ▼
                ┌──────────────────┐
                │ Tem Pulumi?      │
                └──────────────────┘
                   NÃO ↓  ↓ SIM
                       │  │
            ┌──────────┘  └──────────┐
            ▼                         ▼
    ┌──────────────┐        ┌──────────────┐
    │ Instalar     │        │ pulumi login │
    │ Pulumi       │        │ (2 min)      │
    │ (5 min)      │        └──────────────┘
    └──────────────┘                │
            │                        │
            └────────────┬───────────┘
                         ▼
                ┌──────────────────┐
                │ cd infra/pulumi  │
                │ pulumi stack init│
                │ dev              │
                └──────────────────┘
                         │
                         ▼
                ┌──────────────────┐
                │ pulumi preview   │
                └──────────────────┘
                         │
              OK? ────────┼────── ERRO?
              │           │         │
              ▼           ▼         ▼
         ┌─────────┐  Verifique  ┌──────────┐
         │pulumi up│  outputs    │ Troubleshoot │
         │(15 min) │             │ (ver guide)   │
         └─────────┘             └──────────────┘
              │                          │
              └──────────────┬───────────┘
                             ▼
                    ┌──────────────────┐
                    │ Configure GitHub │
                    │ Secrets (5 min)  │
                    └──────────────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │ git push origin  │
                    │ main             │
                    └──────────────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │ GitHub Actions   │
                    │ roda!            │
                    │ (15 min)         │
                    └──────────────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │ ✅ APP RODANDO   │
                    │ NA AWS!          │
                    └──────────────────┘
```

---

## 📚 Ordem de Leitura Recomendada

### 🏃 Para os Apressados (15 min)
1. [QUICK_START.md](./QUICK_START.md)
2. Execute o script
3. Push para GitHub

### 🚶 Para Entender Bem (1-2 horas)
1. [DEPLOYMENT_DOCS_INDEX.md](./DEPLOYMENT_DOCS_INDEX.md) - Overview
2. [AWS_DEPLOYMENT_GUIDE.md](./AWS_DEPLOYMENT_GUIDE.md) - Passo a passo
3. [CI_CD_ANALYSIS.md](./CI_CD_ANALYSIS.md) - Problemas potenciais
4. [DEPLOYMENT_CHECKLIST.md](./DEPLOYMENT_CHECKLIST.md) - Validação

### 🎓 Para Dominar (4-6 horas)
1. Leia tudo acima
2. Execute passo a passo
3. Customize conforme suas necessidades
4. Leia documentação oficial (AWS, Pulumi, Go)
5. Implemente melhorias sugeridas

---

## 🎯 Objetivos Por Fase

### Fase 1: Setup (Objetivo: Ter ferramentas instaladas)
- [ ] AWS CLI instalado
- [ ] Pulumi instalado
- [ ] Docker instalado
- [ ] AWS credentials configuradas
- [ ] Pulumi token obtido

### Fase 2: Infraestrutura (Objetivo: Recursos criados na AWS)
- [ ] Stack Pulumi criado
- [ ] Preview validado
- [ ] Pulumi up executado
- [ ] Outputs capturados
- [ ] Recursos visíveis na AWS Console

### Fase 3: CI/CD (Objetivo: Pipeline automático funcionando)
- [ ] GitHub secrets configurados
- [ ] Workflow atualizado
- [ ] Commit feito
- [ ] Pipeline passou
- [ ] App rodando via Load Balancer

### Fase 4: Validação (Objetivo: Sistema pronto para produção)
- [ ] Health checks passando
- [ ] RDS conectado
- [ ] Logging funcionando
- [ ] Backups configurados
- [ ] Alertas ativados

---

## 🚨 Pontos Críticos

⚠️ **NÃO FAÇA:**
- [ ] Não compartilhe AWS Access Keys
- [ ] Não commite .env com dados reais
- [ ] Não use AdministratorAccess em produção
- [ ] Não deixe Pulumi destruir sem backup
- [ ] Não esqueça de parar recursos em horário de folga

✅ **SEMPRE FAÇA:**
- [ ] Faça backup do Pulumi stack state
- [ ] Use secrets manager para senhas
- [ ] Testar em dev antes de prod
- [ ] Monitorar logs regularmente
- [ ] Revisar custos mensalmente

---

## 💾 Checklist de Segurança

### Antes de Usar em Produção
- [ ] Senhas RDS aleatórias e fortes
- [ ] IAM policy restrita ao mínimo
- [ ] Backup strategy definida
- [ ] Monitoring/Alerting configurado
- [ ] Disaster recovery plan
- [ ] SSL/TLS habilitado (HTTPS)
- [ ] VPC privado configurado
- [ ] WAF nas rules públicas (opcional)

---

## 🔗 Atalhos Rápidos

| Objetivo | Comando |
|----------|---------|
| Ver logs | `pulumi logs` |
| Ver outputs | `pulumi stack output` |
| Deletar tudo | `pulumi destroy` |
| Testar localmente | `docker build -t test .` |
| Ver AWS Resources | `aws ec2 describe-instances` |
| Sync state | `pulumi refresh` |
| Monitor GH Actions | `https://github.com/SEU_REPO/actions` |

---

## 📞 Quando Ficar Perdido

1. **Problema no setup?**
   → Vá para [AWS_DEPLOYMENT_GUIDE.md - Troubleshooting](./AWS_DEPLOYMENT_GUIDE.md#troubleshooting)

2. **Não sabe o que fazer?**
   → Vá para [DEPLOYMENT_CHECKLIST.md](./DEPLOYMENT_CHECKLIST.md)

3. **Quer entender a arquitetura?**
   → Vá para [CI_CD_ANALYSIS.md](./CI_CD_ANALYSIS.md)

4. **Quer começar rápido?**
   → Vá para [QUICK_START.md](./QUICK_START.md)

---

## 🎉 Quando Tudo Estiver Pronto

```bash
# Parabéns! Você conseguiu! 🎉
# Seu app agora:
✅ Roda em AWS (prod-ready)
✅ Tem CI/CD automático
✅ Escala sob demanda
✅ Tem banco de dados gerenciado
✅ Está protegido atrás de LB
✅ Pode fazer deploy com git push

# Próxima missão: Melhorias!
- Auto Scaling
- Blue/Green Deployment
- Multi-region
- Service Mesh
- GitOps
```

---

## 📊 Métrica de Progresso

```
┌────────────────────────────────────────┐
│ Seu Progresso na Jornada AWS          │
├────────────────────────────────────────┤
│ ▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│ 12% - Lendo documentação              │
│                                        │
│ ▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│ 20% - Setup ferramentas                │
│                                        │
│ ▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░░░░ │
│ 30% - AWS Account criada               │
│                                        │
│ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░ │
│ 40% - Pulumi configurado               │
│                                        │
│ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░ │
│ 50% - Primeira execução de deploy      │
│                                        │
│ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░ │
│ 60% - GitHub secrets configurados      │
│                                        │
│ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░ │
│ 70% - CI/CD pipeline rodando           │
│                                        │
│ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░ │
│ 80% - Validações passando              │
│                                        │
│ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░ │
│ 90% - App respondendo na URL           │
│                                        │
│ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ │
│ 100% - PRODUCTION READY! 🚀            │
└────────────────────────────────────────┘
```

---

**Você está no início da jornada! Vamos começar? →** [QUICK_START.md](./QUICK_START.md)

