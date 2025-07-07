#!/bin/bash
set -e

echo "Запуск docker-compose..."
docker-compose up -d

echo "Ожидание запуска базы данных..."
until docker exec ticket_db pg_isready -U postgres; do
  sleep 2
done

echo "Все сервисы запущены!"
docker-compose ps
