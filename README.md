# gestion_cuentas_bancarias
Prueba Técnica Programación | API para la gestión de cuentas bancarias.

# Requerimientos

Tener Instalado Go
https://go.dev/doc/install



# ******************Windows******************

Para ejecutar escribe en la terminal run.bat

Para detener escribe en la terminal stop.bat


# Run Server: Colocar en la terminal CMD en la carpeta base gestion_cuentas_bancarias
go build -o gestioncuentas.exe ./cmd/api
set ENV=development && start /B gestioncuentas.exe

# Stop Server: 
    taskkill /IM gestioncuentas.exe /F

# Eliminar Cache:
    DEL gestioncuentas.exe
    go clean


# ******************Mac OS******************

Usa server.sh

*Dale permisos de ejecucion
chmod +x server.sh

*Dale permisos de ejecucion
./server.sh start    # Iniciar
./server.sh stop     # Detener
./server.sh clean    # Limpiar
./server.sh restart  # Reiniciar

