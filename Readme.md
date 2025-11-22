## Название проекта
Платформа для мастеров ручной работы

## Описание идеи проекта
Создать централизованную платформу, где мастера могут создать свое портфолио, а покупатели — открывать для себя уникальные товары и напрямую связываться с создателями.

## Описание акторов (ролей)

**Мастер** может создать свой магазин, создавать посты о своих работах, в своем портфолио указывать список своих работ(товаров) и свои соц сети для связи с покупателем.

**Пользователь**, может смотреть посты, по категориям и искать мастеров.

## Стек
Docker Compose, Nginx, gRPC, REST API, Go, JWT, Swagger, Postgres


## Use-Case - диаграмма
![Use-Case](img/usecase_craftPlace.png)

## ER-диаграмма сущностей
![ER-диаграмма](img/ER_craftPlace.png)

## 10. Формализация ключевых бизнес-процессов (BPMN-нотация).
![BPMN](img/bpmn_craftPlace.png)

## Верхнеуровневое разбиение на компоненты
![BPMN](img/components_craftPlace.png)

## Черновик интерфейса
![interface](img/interface_craftPlace.jpg)

## Документация (Swagger)
[swagger.yaml](./docs/swagger.yaml)

[http://localhost/swagger/index.html](http://localhost/swagger/index.html)

Loki: http://loki_craftplace:3100




# Nginx Configuration Paths for README.md

## Основные маршруты приложения

### 📍 Статические файлы и основной интерфейс
- **`/`** - Главная страница (статический контент)
- **`/status`** - Страница статуса сервиса
- **`/documentation`** - Документация проекта
- **`/managment`** - Панель управления

### 📚 Документация API
- **`/swagger/`** → Primary Backend
  - Swagger документация основного API
- **`/mirror/swagger/`** → Mirror Backend
  - Swagger документация mirror API


#### Администрирование БД
- **`/admin/`** → pgAdmin
  - Веб-интерфейс для управления PostgreSQL
### 🔄 Mirror API
- **`/mirror/api/v1/`** → Mirror Backend
  - Тестовое/зеркальное API для разработки





































  
