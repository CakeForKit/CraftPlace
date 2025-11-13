#!/bin/bash

ab -n 1000 -c 1 "http://localhost:80/api/v1/categories?page=1&size=20"
