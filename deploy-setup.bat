@echo off
REM Deploy Setup Script para Windows
REM Uso: deploy-setup.bat

setlocal enabledelayedexpansion

echo ================================================
echo   Deploy Setup - AWS Deployment
echo ================================================
echo.

REM Check 1: AWS CLI
echo [1/5] Verificando AWS CLI...
where aws >nul 2>nul
if errorlevel 1 (
    echo X AWS CLI nao instalado
    echo   Instale com: pip install awscli
    exit /b 1
)
echo ^+ AWS CLI encontrado
for /f "tokens=*" %%A in ('aws --version') do set AWS_VER=%%A
echo   %AWS_VER%

REM Check 2: Pulumi
echo.
echo [2/5] Verificando Pulumi...
where pulumi >nul 2>nul
if errorlevel 1 (
    echo X Pulumi nao instalado
    echo   Windows: choco install pulumi
    exit /b 1
)
echo ^+ Pulumi encontrado
for /f "tokens=*" %%A in ('pulumi version') do set PULUMI_VER=%%A
echo   %PULUMI_VER%

REM Check 3: Docker
echo.
echo [3/5] Verificando Docker...
where docker >nul 2>nul
if errorlevel 1 (
    echo X Docker nao instalado
    echo   Instale em: https://www.docker.com/products/docker-desktop
    exit /b 1
)
echo ^+ Docker encontrado
for /f "tokens=*" %%A in ('docker --version') do set DOCKER_VER=%%A
echo   %DOCKER_VER%

REM Check 4: Go
echo.
echo [4/5] Verificando Go...
where go >nul 2>nul
if errorlevel 1 (
    echo X Go nao instalado
    exit /b 1
)
echo ^+ Go encontrado
for /f "tokens=3" %%A in ('go version') do set GO_VER=%%A
echo   %GO_VER%

REM Check 5: AWS Credentials
echo.
echo [5/5] Verificando AWS Credentials...
aws sts get-caller-identity >nul 2>nul
if errorlevel 1 (
    echo X AWS Credentials nao configurados
    echo   Execute: aws configure
    exit /b 1
)
echo ^+ Credenciais AWS configuradas
for /f "tokens=*" %%A in ('aws sts get-caller-identity --query Account --output text') do set ACCOUNT_ID=%%A
echo   AWS Account: %ACCOUNT_ID%

REM Tudo certo!
echo.
echo ================================================
echo + Todas as dependencias estao OK!
echo ================================================
echo.

echo PROXIMOS PASSOS:
echo.
echo 1. Preparar Pulumi:
echo    cd infra\pulumi
echo    pulumi login
echo    pulumi stack init dev
echo    pulumi config set aws:region us-east-1
echo.
echo 2. Preview do deploy:
echo    pulumi preview
echo.
echo 3. Fazer deploy:
echo    pulumi up
echo.
echo 4. Configurar GitHub Secrets
echo    https://github.com/SEU_USUARIO/SEU_REPO/settings/secrets/actions
echo.
echo 5. Fazer push para main
echo.
echo Para mais detalhes: AWS_DEPLOYMENT_GUIDE.md
echo.
pause
