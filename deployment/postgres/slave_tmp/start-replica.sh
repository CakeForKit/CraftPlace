#!/bin/bash

pg_basebackup -h postgres_craftplace -U repl_user -D /var/lib/postgresql/data -P --wal-method=stream --slot=replication_slot

cp /scripts/recovery.conf /var/lib/postgresql/data/

# Запуск реплики
exec postgres -c hot_standby=on