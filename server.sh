#!/bin/bash

case "$1" in
  start)
    echo "Compilando y ejecutando el servidor..."
    go build -o gestioncuentas ./cmd/api
    ENV=development nohup ./gestioncuentas > output.log 2>&1 &
    echo "Servidor iniciado!"
    ;;
    
  stop)
    echo "Deteniendo el servidor..."
    pkill -f gestioncuentas
    echo "Servidor detenido!"
    ;;
    
  clean)
    echo "Eliminando binarios y limpiando caché..."
    rm -f gestioncuentas
    go clean
    echo "Limpieza completada!"
    ;;
    
  restart)
    echo "Reiniciando el servidor..."
    pkill -f gestioncuentas
    go build -o gestioncuentas ./cmd/api
    ENV=development nohup ./gestioncuentas > output.log 2>&1 &
    echo "Servidor reiniciado!"
    ;;
    
  *)
    echo "Uso: ./server.sh {start|stop|clean|restart}"
    exit 1
    ;;
esac
