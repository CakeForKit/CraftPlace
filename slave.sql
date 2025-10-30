-- Active: 1761818804170@@127.0.0.1@5433@craftplace

SELECT TABLE_NAME
FROM INFORMATION_SCHEMA.TABLES;

select * from pg_replication_slots;

SELECT * FROM shops;
SELECT * FROM products;
SELECT * FROM posts;
SELECT * FROM users;
SELECT * FROM categories;

UPDATE shops 
SET title = 'Лавка "искусства" 1', 
    update_time = CURRENT_TIMESTAMP 
WHERE id = 'e855444e-26a4-47b4-b7fe-f15ff7cb552e';