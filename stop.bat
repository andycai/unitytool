@echo off
set SERVICE_NAME=unitool_serve_linux.exe
if "%1" neq "" set SERVICE_NAME=%1
for /f "tokens=2" %%i in ('tasklist ^| findstr /I "%SERVICE_NAME%"') do taskkill /PID %%i /F