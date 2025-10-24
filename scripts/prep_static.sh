#!/bin/bash

pwd 
pandoc ./Readme.md -f markdown -t html -o ./nginx/static/Readme.html

./scripts/download-swagger-static.sh

mkdir -p ./nginx/static/img/
cp ./img/* ./nginx/static/img/
