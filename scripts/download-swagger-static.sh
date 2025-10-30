#!/bin/bash

# Настройки
SWAGGER_BASE_URL="http://localhost:8080/mirror/swagger"
TARGET_DIR="./nginx/static/swagger"

# Создаем целевую папку
mkdir -p "$TARGET_DIR"
# cd "$TARGET_DIR"

echo "Downloading custom Swagger files from $SWAGGER_BASE_URL..."

# Список файлов для загрузки (из HTML)
FILES=(
    "index.html"
    "swagger-ui.css"
    "favicon-32x32.png" 
    "favicon-16x16.png"
    "index.css"
    "swagger-ui-bundle.js"
    "swagger-ui-standalone-preset.js"
    "swagger-initializer.js"
)

for file in "${FILES[@]}"; do
    curl -s -f -o "$TARGET_DIR/$file" "http://localhost:8080/mirror/swagger/$file"
done

cp ./docs/* "$TARGET_DIR"
