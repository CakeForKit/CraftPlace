#!/bin/bash

# Настройки
SWAGGER_BASE_URL="http://localhost:8080/swagger"
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
    # echo "Executing: curl -s -f -o \"$TARGET_DIR/$file\" \"$SWAGGER_BASE_URL/$file\""
    curl -s -f -o "$TARGET_DIR/$file" "http://localhost:8080/swagger/$file"
done

cp ./docs/* "$TARGET_DIR"


# # Загружаем каждый файл
# for file in "${FILES[@]}"; do
#     echo -n "Downloading $file... "
    
#     # Используем curl с опциями для лучшей диагностики
#     if curl -s -f -o "$file" "$SWAGGER_BASE_URL/$file"; then
#         # Проверяем что файл не пустой и не содержит HTML ошибку
#         if [ -s "$file" ]; then
#             first_line=$(head -1 "$file" | tr -d '\n\r')
#             if [[ "$first_line" == "<!DOCTYPE html>"* ]] || [[ "$first_line" == "<html"* ]]; then
#                 echo "✗ FAILED (contains HTML, not actual file)"
#                 rm -f "$file"
#             else
#                 file_size=$(stat -c%s "$file" 2>/dev/null || stat -f%z "$file" 2>/dev/null)
#                 echo "✓ SUCCESS ($file_size bytes)"
#             fi
#         else
#             echo "✗ FAILED (empty file)"
#             rm -f "$file"
#         fi
#     else
#         echo "✗ FAILED (download error)"
#     fi
# done

# # Также пробуем загрузить doc.json
# echo "Downloading doc.json..."
# if curl -f -s "$SWAGGER_BASE_URL/doc.json" -o "doc.json"; then
#     echo "✓ Downloaded doc.json"
# else
#     echo "✗ Failed to download doc.json"
# fi

# # Проверяем что скачалось
# echo ""
# echo "Downloaded files:"
# ls -la

# # Проверяем размеры файлов
# echo ""
# echo "File sizes:"
# for file in *; do
#     if [ -f "$file" ]; then
#         size=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null)
#         echo "  $file: $size bytes"
#     fi
# done

# echo "Done!"
