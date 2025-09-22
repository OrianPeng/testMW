@echo off
echo Starting AI Assistant Frontend...
echo.

REM 检查Node.js是否安装
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Node.js is not installed or not in PATH
    echo Please install Node.js from https://nodejs.org/
    pause
    exit /b 1
)

REM 检查package.json是否存在
if not exist package.json (
    echo Error: package.json not found
    echo Please run this script from the frontend directory
    pause
    exit /b 1
)

REM 安装依赖（如果需要）
if not exist node_modules (
    echo Installing dependencies...
    npm install
    if %errorlevel% neq 0 (
        echo Error: Failed to install dependencies
        pause
        exit /b 1
    )
)

REM 启动开发服务器
echo Starting development server...
echo.
echo AI Assistant will be available at: http://localhost:3000/chat
echo Test iframe page: http://localhost:3000/test-iframe.html
echo.
echo Press Ctrl+C to stop the server
echo.

npm run dev

