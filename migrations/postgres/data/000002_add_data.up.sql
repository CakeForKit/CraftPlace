-- Active: 1760105453579@@127.0.0.1@5432@craftplace


-- Active: 1744740356603@@127.0.0.1@5432@artworks

-- Функция для генерации случайного логина
CREATE OR REPLACE FUNCTION random_login() RETURNS VARCHAR(50) AS $$
DECLARE
    prefixes VARCHAR[] := ARRAY['user', 'client', 'buyer', 'seller', 'shopowner', 'artlover', 'creator', 'designer', 'maker', 'craftsman'];
    suffixes VARCHAR[] := ARRAY['123', 'pro', 'max', 'top', 'best', 'cool', 'art', 'shop', 'store', 'online'];
BEGIN
    RETURN prefixes[1 + floor(random() * array_length(prefixes, 1))] || '_' || 
           suffixes[1 + floor(random() * array_length(suffixes, 1))] || 
           floor(random() * 1000)::INT;
END;
$$ LANGUAGE plpgsql;

-- Функция для генерации хеша пароля (имитация)
CREATE OR REPLACE FUNCTION random_password_hash() RETURNS VARCHAR(255) AS $$
BEGIN
    RETURN 'hashed_password_' || substr(md5(random()::text), 0, 20);
END;
$$ LANGUAGE plpgsql;

-- Функция для генерации названия магазина
CREATE OR REPLACE FUNCTION random_shop_title() RETURNS VARCHAR(255) AS $$
DECLARE
    prefixes VARCHAR[] := ARRAY['Магазин', 'Бутик', 'Студия', 'Галерея', 'Мастерская', 'Ателье', 'Лавка', 'Салон', 'Компания', 'Торговая точка'];
    suffixes VARCHAR[] := ARRAY['уникальных товаров', 'рукоделия', 'творчества', 'искусства', 'дизайна', 'красоты', 'стиля', 'моды', 'хендмейда', 'авторских работ'];
BEGIN
    RETURN prefixes[1 + floor(random() * array_length(prefixes, 1))] || ' "' || 
           suffixes[1 + floor(random() * array_length(suffixes, 1))] || '"';
END;
$$ LANGUAGE plpgsql;

-- Функция для генерации описания магазина
CREATE OR REPLACE FUNCTION random_shop_description() RETURNS VARCHAR(500) AS $$
DECLARE
    descriptions VARCHAR[] := ARRAY[
        'Мы предлагаем уникальные товары ручной работы от талантливых мастеров.',
        'Широкий ассортимент качественных продуктов для творческих людей.',
        'Авторские работы и эксклюзивные товары только в нашем магазине.',
        'Лучшие товары для вашего вдохновения и творчества.',
        'Мы заботимся о качестве и уникальности каждого продукта.',
        'Магазин для тех, кто ценит ручную работу и индивидуальность.',
        'Только проверенные мастера и качественные материалы.',
        'Ваш надежный партнер в мире творчества и рукоделия.',
        'Уникальные решения для вашего бизнеса и хобби.',
        'Мы делаем мир красивее, один товар за разным.'
    ];
BEGIN
    RETURN descriptions[1 + floor(random() * array_length(descriptions, 1))];
END;
$$ LANGUAGE plpgsql;

-- Функция для генерации названия продукта
CREATE OR REPLACE FUNCTION random_product_title() RETURNS VARCHAR(255) AS $$
DECLARE
    prefixes VARCHAR[] := ARRAY['Уникальный', 'Эксклюзивный', 'Авторский', 'Ручной работы', 'Дизайнерский', 'Стильный', 'Креативный', 'Элегантный', 'Современный', 'Винтажный'];
    products VARCHAR[] := ARRAY['браслет', 'портрет', 'сумка', 'свитер', 'горшок', 'светильник', 'альбом', 'открытка', 'украшение', 'игрушка', 'картина', 'шкатулка', 'зеркало', 'подсвечник', 'ваза'];
BEGIN
    RETURN prefixes[1 + floor(random() * array_length(prefixes, 1))] || ' ' || 
           products[1 + floor(random() * array_length(products, 1))];
END;
$$ LANGUAGE plpgsql;

-- Функция для генерации описания продукта
CREATE OR REPLACE FUNCTION random_product_description() RETURNS VARCHAR(500) AS $$
DECLARE
    descriptions VARCHAR[] := ARRAY[
        'Этот товар создан с любовью и вниманием к деталям.',
        'Качественные материалы и тщательная обработка.',
        'Идеальный подарок для близких и друзей.',
        'Уникальный дизайн, который выделит вас из толпы.',
        'Сочетание традиционных техник и современных тенденций.',
        'Экологически чистые материалы и безопасное производство.',
        'Продукт прошел строгий контроль качества.',
        'Ограниченная серия - только для настоящих ценителей.',
        'Идеально подходит для подарка на любой праздник.',
        'Сохранит ваши воспоминания на долгие годы.'
    ];
BEGIN
    RETURN descriptions[1 + floor(random() * array_length(descriptions, 1))];
END;
$$ LANGUAGE plpgsql;

-- Функция для генерации описания поста
CREATE OR REPLACE FUNCTION random_post_description() RETURNS VARCHAR(500) AS $$
DECLARE
    descriptions VARCHAR[] := ARRAY[
        'Новые поступления в нашем магазине! Не пропустите уникальные товары.',
        'Специальное предложение для наших постоянных клиентов.',
        'Рассказываем о процессе создания наших продуктов.',
        'Идеи для творчества и вдохновения от нашей команды.',
        'Отзывы довольных клиентов и их покупки.',
        'Мастер-класс по использованию наших товаров.',
        'Новая коллекция уже в продаже! Успейте приобрести.',
        'Акция недели: скидки на популярные товары.',
        'Интервью с нашими мастерами и дизайнерами.',
        'Полезные советы по уходу за изделиями ручной работы.'
    ];
BEGIN
    RETURN descriptions[1 + floor(random() * array_length(descriptions, 1))];
END;
$$ LANGUAGE plpgsql;

-- Функция для генерации названия категории
CREATE OR REPLACE FUNCTION random_category_title() RETURNS VARCHAR(255) AS $$
DECLARE
    categories VARCHAR[] := ARRAY['Бижутерия', 'Одежда', 'Аксессуары', 'Декор', 'Канцелярия', 'Игрушки', 'Посуда', 'Текстиль', 'Сувениры', 'Материалы для творчества'];
BEGIN
    RETURN categories[1 + floor(random() * array_length(categories, 1))];
END;
$$ LANGUAGE plpgsql;

-- Функция для генерации описания категории
CREATE OR REPLACE FUNCTION random_category_description() RETURNS VARCHAR(500) AS $$
DECLARE
    descriptions VARCHAR[] := ARRAY[
        'Разнообразные товары в этой категории для вашего творчества.',
        'Лучшие продукты от проверенных мастеров и дизайнеров.',
        'Широкий выбор качественных товаров по доступным ценам.',
        'Уникальные решения для вашего дома и подарков.',
        'Товары, которые вдохновят на новые творческие идеи.',
        'Эксклюзивные продукты, которые вы не найдете больше нигде.',
        'Сочетание традиционного мастерства и современных тенденций.',
        'Идеальные товары для начинающих и опытных мастеров.',
        'Вся продукция проходит строгий отбор по качеству.',
        'Категория, которая постоянно пополняется новинками.'
    ];
BEGIN
    RETURN descriptions[1 + floor(random() * array_length(descriptions, 1))];
END;
$$ LANGUAGE plpgsql;

-- Заполнение таблицы users (30 пользователей)
DO $$
DECLARE
    i INTEGER;
BEGIN
    FOR i IN 1..30 LOOP
        INSERT INTO users (id, login, hashed_password)
        VALUES (
            gen_random_uuid(),
            random_login(),
            random_password_hash()
        );
    END LOOP;
    RAISE NOTICE 'Добавлено % пользователей', (SELECT COUNT(*) FROM users);
END $$;

-- Заполнение таблицы shops (15 магазинов)
DO $$
DECLARE
    i INTEGER;
    user_rec RECORD;
BEGIN
    FOR i IN 1..15 LOOP
        -- Выбираем случайного пользователя
        SELECT id INTO user_rec FROM users ORDER BY random() LIMIT 1;
        
        INSERT INTO shops (id, title, description, user_id, update_time)
        VALUES (
            gen_random_uuid(),
            random_shop_title(),
            random_shop_description(),
            user_rec.id,
            NOW() - (random() * 365 || ' days')::INTERVAL
        );
    END LOOP;
    RAISE NOTICE 'Добавлено % магазинов', (SELECT COUNT(*) FROM shops);
END $$;

-- Заполнение таблицы posts (50 постов)
DO $$
DECLARE
    i INTEGER;
    shop_rec RECORD;
BEGIN
    FOR i IN 1..50 LOOP
        -- Выбираем случайный магазин
        SELECT id INTO shop_rec FROM shops ORDER BY random() LIMIT 1;
        
        INSERT INTO posts (id, description, publication_time, shop_id)
        VALUES (
            gen_random_uuid(),
            random_post_description(),
            NOW() - (random() * 30 || ' days')::INTERVAL,
            shop_rec.id
        );
    END LOOP;
    RAISE NOTICE 'Добавлено % постов', (SELECT COUNT(*) FROM posts);
END $$;

-- Заполнение таблицы categories (30 категорий)
DO $$
DECLARE
    i INTEGER;
BEGIN
    FOR i IN 1..30 LOOP
        INSERT INTO categories (id, title, description)
        VALUES (
            gen_random_uuid(),
            random_category_title(),
            random_category_description()
        );
    END LOOP;
    RAISE NOTICE 'Добавлено % категорий', (SELECT COUNT(*) FROM categories);
END $$;

-- Заполнение таблицы products (40 продуктов)
DO $$
DECLARE
    i INTEGER;
    shop_rec RECORD;
BEGIN
    FOR i IN 1..40 LOOP
        -- Выбираем случайный магазин
        SELECT id INTO shop_rec FROM shops ORDER BY random() LIMIT 1;
        
        INSERT INTO products (id, title, description, cost, shop_id, update_time)
        VALUES (
            gen_random_uuid(),
            random_product_title(),
            random_product_description(),
            100 + floor(random() * 4900)::INT, -- стоимость от 100 до 5000
            shop_rec.id,
            NOW() - (random() * 60 || ' days')::INTERVAL
        );
    END LOOP;
    RAISE NOTICE 'Добавлено % продуктов', (SELECT COUNT(*) FROM products);
END $$;

-- Заполнение таблицы product_category (связи продуктов с категориями)
DO $$
DECLARE
    product_rec RECORD;
    category_rec RECORD;
    categories_count INTEGER;
    i INTEGER;
BEGIN
    -- Для каждого продукта добавляем 1-5 случайные категории
    FOR product_rec IN SELECT id FROM products LOOP
        categories_count := 1 + floor(random() * 5)::INT; -- от 1 до 5 категорий на продукт
        
        FOR i IN 1..categories_count LOOP
            -- Выбираем случайную категорию
            SELECT id INTO category_rec FROM categories ORDER BY random() LIMIT 1;
            
            INSERT INTO product_category (product_id, category_id)
            VALUES (product_rec.id, category_rec.id)
            ON CONFLICT (product_id, category_id) DO NOTHING;
        END LOOP;
    END LOOP;
    
    RAISE NOTICE 'Добавлено связей продукт-категория: %', (SELECT COUNT(*) FROM product_category);
END $$;

-- Вывод итоговой статистики
DO $$
BEGIN
    RAISE NOTICE '=== ИТОГОВАЯ СТАТИСТИКА ===';
    RAISE NOTICE 'Пользователей: %', (SELECT COUNT(*) FROM users);
    RAISE NOTICE 'Магазинов: %', (SELECT COUNT(*) FROM shops);
    RAISE NOTICE 'Постов: %', (SELECT COUNT(*) FROM posts);
    RAISE NOTICE 'Продуктов: %', (SELECT COUNT(*) FROM products);
    RAISE NOTICE 'Категорий: %', (SELECT COUNT(*) FROM categories);
    RAISE NOTICE 'Связей продукт-категория: %', (SELECT COUNT(*) FROM product_category);
END $$;