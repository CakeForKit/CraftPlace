-- Active: 1760105453579@@127.0.0.1@5432@craftplace

SELECT TABLE_NAME
FROM INFORMATION_SCHEMA.TABLES;

select * from pg_replication_slots;

SELECT * FROM pg_stat_replication;

SELECT slot_name, slot_type, active, wal_status FROM pg_replication_slots;

SELECT rolname, rolreplication FROM pg_roles WHERE rolname = 'repl_user';
select * from pg_roles;

SELECT * FROM shops;
SELECT * FROM products;
SELECT * FROM posts;
SELECT * FROM users WHERE login = 'ulogin';
SELECT * FROM categories;

UPDATE shops 
SET title = 'Лавка "искусства" 1', 
    update_time = CURRENT_TIMESTAMP 
WHERE id = 'e855444e-26a4-47b4-b7fe-f15ff7cb552e';
