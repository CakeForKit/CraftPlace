project/
├── proto/                          # Protobuf схемы
│   ├── user_service.proto
│   ├── order_service.proto
│   └── common/                     # Общие сообщения
│       ├── common.proto
│       └── types.proto
├── internal/                       # Внутренние библиотеки
│   ├── config/                     # Конфигурация
│   ├── database/                   # Подключение к БД
│   └── middleware/                 # gRPC middleware
├── pkg/                           # Переиспользуемые пакеты
│   ├── grpcclient/                # gRPC клиенты
│   ├── logger/                    # Логирование
│   └── utils/                     # Утилиты
├── user-service/                  # Сервис пользователей
│   ├── cmd/
│   │   └── main.go
│   ├── internal/                  # Логика сервиса
│   │   ├── handler/               # gRPC handlers
│   │   ├── service/               # Бизнес-логика
│   │   ├── repository/            # Работа с данными
│   │   └── models/                # Модели данных
│   └── Dockerfile
├── order-service/                 # Сервис заказов
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   └── models/
│   └── Dockerfile
└── api-gateway/                   # API Gateway
    ├── cmd/
    │   └── main.go
    └── internal/
        └── http/


        