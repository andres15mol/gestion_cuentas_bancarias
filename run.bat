@echo off
echo Iniciando el backend...
go build -o gestioncuentas.exe ./cmd/api
set ENV=development && start /B gestioncuentas.exe
echo Backend iniciado!
