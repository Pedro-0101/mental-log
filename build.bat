@echo off
echo Building Mental Log...
go build -o mental-log.exe ./cmd/main.go
if %errorlevel% neq 0 (
    echo Build failed!
    pause
    exit /b %errorlevel%
)
echo Build success!
echo Running Mental Log...
./mental-log.exe
pause
