SCRIPTS := ./scripts
DC_CI := ./deployment/docker-compose.ci.yml
DC_DEV := ./deployment/docker-compose.dev.yml
TEST_DB_ENV := ./configs/test_db.env

.PHONY: run_app
run_app:
# --no-cache
	docker compose -v -f $(DC_DEV) build --progress=plain app_craftplace
	docker compose -v -f $(DC_DEV) up  app_craftplace

.PHONY: down_app
down_app:
	docker compose -f $(DC_DEV) down -v app_craftplace

.PHONY: swagger
swagger:
	swag init -g ./cmd/dev/main.go --output ./docs


# ---- Allure -----
ALLURE_OUTPUT_PATH := $(shell pwd)
ALLURE_RESULTS_DIR := $(shell pwd)/allure-results
ALLURE_REPORT_DIR := $(shell pwd)/allure-report
export ALLURE_RESULTS_DIR
export ALLURE_OUTPUT_PATH

.PHONY: allure
allure: unit_test report_allure open_allure

.PHONY: unit_test
unit_test : clear_allure
	$(SCRIPTS)/unit_tests.sh

# .PHONY: integration_test
# integration_test: clear_allure
# 	$(SCRIPTS)/integration_tests.sh

.PHONY: report_allure
report_allure:
	mkdir -p $(ALLURE_REPORT_DIR)/history
	cp -r $(ALLURE_REPORT_DIR)/history $(ALLURE_RESULTS_DIR)
	allure generate $(ALLURE_RESULTS_DIR) -o $(ALLURE_REPORT_DIR) --clean

.PHONY: clear_allure
clear_allure:
	rm -rf $(ALLURE_RESULTS_DIR)

.PHONY: open_allure
open_allure:
	allure open $(ALLURE_REPORT_DIR)
	

.PHONY: test_run_app
test_run_app:
# --no-cache
	docker compose -v -f $(DC_CI) --env-file $(TEST_DB_ENV) build --progress=plain test-runner
	docker compose -v -f $(DC_CI) --env-file $(TEST_DB_ENV) up  test-runner

.PHONY: test_down_app
test_down_app:
	docker compose -f $(DC_CI) down -v test-runner


.PHONY: test_run_serv
test_run_serv:
	docker compose -f $(DC_CI) --env-file $(TEST_DB_ENV) up -d postgres migrator redis_artworks

.PHONY: test_down_serv
test_down_serv:
	docker compose -f $(DC_CI) down -v postgres migrator redis_artworks

.PHONY: test_build
test_build:
	docker compose -f $(DC_CI) --env-file $(TEST_DB_ENV) build


.PHONY: clear_docker
clear_docker:
# Остановите все контейнеры
	docker-compose -f ./deployment/docker-compose.ci.yml down
# Удалите старые образы
	docker rmi deployment-test-runner
# Очистите builder кеш
	docker builder prune -f
# Удалите все старые версии
	docker rmi deployment-test-runner:latest
# Полная очистка
	docker system prune -a -f



