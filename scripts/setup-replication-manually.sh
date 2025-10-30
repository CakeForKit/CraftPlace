#!/bin/bash

echo "Настройка репликации PostgreSQL..."

echo "Ожидаем запуск PostgreSQL..."
until docker exec postgres_craftplace pg_isready -U puser -d craftplace; do
  sleep 2
done

echo "PostgreSQL готов, настраиваем репликацию..."

docker exec postgres_craftplace psql -U puser -d craftplace -c "
CREATE ROLE repl_user WITH REPLICATION LOGIN PASSWORD 'repl_password';
"

docker exec postgres_craftplace psql -U puser -d craftplace -c "
SELECT pg_create_physical_replication_slot('replication_slot');
"

docker exec postgres_craftplace sh -c "
echo 'host replication all 0.0.0.0/0 md5' >> /var/lib/postgresql/data/pg_hba.conf
echo 'host all all 0.0.0.0/0 md5' >> /var/lib/postgresql/data/pg_hba.conf
"

docker exec postgres_craftplace psql -U puser -d craftplace -c "
SELECT pg_reload_conf();
"

echo "Репликация настроена!"
echo "Пользователь: repl_user"
echo "Пароль: repl_password" 
echo "Слот: replication_slot"