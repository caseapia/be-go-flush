@echo off
REM ============================
REM Build Go project for Linux
REM ============================

REM Настройки
set SERVICE_NAME=be-go-flush
set GOARCH=amd64
set GOOS=linux
set PROJECT_PATH=D:\projects\be-go-flush
set BUILD_PATH=%PROJECT_PATH%\build

REM Создаем папку build если нет
if not exist "%BUILD_PATH%" mkdir "%BUILD_PATH%"

echo Starting Go build for service %SERVICE_NAME%...
cd /d "%PROJECT_PATH%"

REM Основной билд
go build -o "%BUILD_PATH%\%SERVICE_NAME%" ./cmd/app
if errorlevel 1 (
    echo Error during Go build process
    pause
    exit /b 1
)

echo Build completed successfully.
echo Files are in %BUILD_PATH%
pause
