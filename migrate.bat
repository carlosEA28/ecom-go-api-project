@echo off
REM Script para executar migrations localmente ou em RDS (Windows)
REM Uso: migrate.bat [dev|staging|prod] [up|down|status|redo]

setlocal enabledelayedexpansion

set ENVIRONMENT=%1
set ACTION=%2

if "%ENVIRONMENT%"=="" (
    echo Error: Environment not specified
    echo Usage: migrate.bat [dev^|staging^|prod] [up^|down^|status^|redo]
    echo.
    echo Examples:
    echo   migrate.bat dev                  - Apply migrations
    echo   migrate.bat dev up               - Apply all pending
    echo   migrate.bat dev down             - Revert last
    echo   migrate.bat dev status           - Check status
    echo   migrate.bat dev redo             - Revert and reapply
    exit /b 1
)

if "%ACTION%"=="" set ACTION=up

REM Validar ambiente
if not "%ENVIRONMENT%"=="dev" if not "%ENVIRONMENT%"=="staging" if not "%ENVIRONMENT%"=="prod" (
    echo Error: Invalid environment "%ENVIRONMENT%"
    echo Use: dev, staging, or prod
    exit /b 1
)

REM Validar ação
if not "%ACTION%"=="up" if not "%ACTION%"=="down" if not "%ACTION%"=="status" if not "%ACTION%"=="redo" (
    echo Error: Invalid action "%ACTION%"
    echo Use: up, down, status, or redo
    exit /b 1
)

echo [*] Environment: %ENVIRONMENT%
echo [*] Action: %ACTION%
echo.

REM Verificar se goose está instalado
where goose >nul 2>nul
if %errorlevel% neq 0 (
    echo [!] Goose not installed. Installing...
    go install github.com/pressly/goose/v3/cmd/goose@latest
    if %errorlevel% neq 0 (
        echo [x] Failed to install Goose
        exit /b 1
    )
)

if "%ENVIRONMENT%"=="dev" (
    echo [*] Using local PostgreSQL (docker-compose)
    
    REM Verificar se docker está rodando
    docker ps | findstr postgres >nul 2>nul
    if %errorlevel% neq 0 (
        echo [!] PostgreSQL not running in Docker
        echo [*] Starting docker-compose...
        docker-compose up -d postgres
        timeout /t 10 /nobreak
    )
    
    set DB_HOST=localhost
    set DB_PORT=5432
    set DB_USER=postgres
    set DB_PASSWORD=postgres
    set DB_NAME=ecom
    set DB_SSLMODE=disable
) else (
    echo [*] Getting credentials from Pulumi for %ENVIRONMENT%...
    
    REM Verificar se pulumi está instalado
    where pulumi >nul 2>nul
    if %errorlevel% neq 0 (
        echo [x] Pulumi not installed
        exit /b 1
    )
    
    cd infra\pulumi
    
    REM Selecionar stack
    pulumi stack select %ENVIRONMENT% >nul 2>&1
    if %errorlevel% neq 0 (
        echo [x] Stack '%ENVIRONMENT%' does not exist
        exit /b 1
    )
    
    REM Obter valores
    for /f "tokens=*" %%i in ('pulumi stack output rdsEndpoint 2^>nul') do set DB_HOST=%%i
    for /f "tokens=*" %%i in ('pulumi stack output rdsPort 2^>nul') do set DB_PORT=%%i
    for /f "tokens=*" %%i in ('pulumi stack output rdsUsername 2^>nul') do set DB_USER=%%i
    for /f "tokens=*" %%i in ('pulumi stack output rdsDatabase 2^>nul') do set DB_NAME=%%i
    
    set DB_SSLMODE=require
    
    if "%RDS_PASSWORD%"=="" (
        set /p DB_PASSWORD=[*] Enter RDS password: 
    ) else (
        set DB_PASSWORD=%RDS_PASSWORD%
    )
    
    cd ..\..
)

REM Construir connection string
set GOOSE_DBSTRING=host=!DB_HOST! port=!DB_PORT! user=!DB_USER! password=!DB_PASSWORD! dbname=!DB_NAME! sslmode=!DB_SSLMODE!

echo [+] Connecting to !DB_HOST!:!DB_PORT!/!DB_NAME!
echo.

REM Executar migrations
cd internal\adapters\postgresql

if "%ACTION%"=="up" (
    echo [+] Applying pending migrations...
    goose postgres "!GOOSE_DBSTRING!" up
    if %errorlevel% neq 0 (
        echo [x] Error applying migrations
        exit /b 1
    )
    echo [✓] Migrations applied successfully
) else if "%ACTION%"=="down" (
    echo [+] Reverting last migration...
    goose postgres "!GOOSE_DBSTRING!" down
    if %errorlevel% neq 0 (
        echo [x] Error reverting migration
        exit /b 1
    )
    echo [✓] Migration reverted successfully
) else if "%ACTION%"=="status" (
    echo [*] Migration status:
    echo.
    goose postgres "!GOOSE_DBSTRING!" status
) else if "%ACTION%"=="redo" (
    echo [+] Reverting and reapplying last migration...
    goose postgres "!GOOSE_DBSTRING!" redo
    if %errorlevel% neq 0 (
        echo [x] Error doing redo
        exit /b 1
    )
    echo [✓] Migration redo completed
)

echo.
echo [✓] Operation completed!
echo.
echo [*] Summary:
echo     Environment: %ENVIRONMENT%
echo     Host: !DB_HOST!:!DB_PORT!
echo     Database: !DB_NAME!
echo     Action: %ACTION%
echo.

echo [*] Current migration status:
goose postgres "!GOOSE_DBSTRING!" status

cd ..\..\..

endlocal
