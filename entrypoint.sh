#!/bin/sh
set -e

# Espera o Postgres ficar pronto
until psql "$DATABASE_URL" -c '\q' 2>/dev/null; do
  echo "Aguardando Postgres..."
  sleep 2
done

echo "Rodando migrations..."
psql "$DATABASE_URL" -f /app/internal/db/migrations.sql

echo "Iniciando aplicação..."
exec "$@"