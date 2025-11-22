SCRIPTS := ./scripts
DC_CI := ./craftplace-deployment/docker-compose.ci.yml
DC_DEV := ./craftplace-deployment/docker-compose.dev.yml
TEST_DB_ENV := ./configs/test_db.env
DB_ENV := ./configs/db_config.env

.PHONY: gen
gen:
	protoc -I proto ./proto/auth_user.proto --go_out=./ --go-grpc_out=./

.PHONY: prep
prep: run_db run_pgadmin serv
	
.PHONY: serv
serv: app_run auth_run searcher_run run_mirror run_loki_graf run_ng 

.PHONY: app_run
app_run:
	docker compose -v -f $(DC_DEV) up --build app_craftplace app_craftplace_replica1 app_craftplace_replica2 -d


.PHONY: auth_run
auth_run:
	docker compose -v -f $(DC_DEV) up --build auth_craftplace auth_craftplace_replica1 auth_craftplace_replica2 -d

.PHONY: searcher_run
searcher_run:
	docker compose -v -f $(DC_DEV) up --build searcher_craftplace searcher_craftplace_replica1 searcher_craftplace_replica2 -d

.PHONY: run_loki_graf
run_loki_graf:
	docker compose -v -f $(DC_DEV) up --build loki_craftplace grafana_craftplace promtail_craftplace -d

.PHONY: down_loki_graf
down_loki_graf:
	docker compose -f $(DC_DEV) down -v loki_craftplace grafana_craftplace promtail_craftplace

.PHONY: run_app
run_app:
# --no-cache	--progress=plain
	docker compose -v -f $(DC_DEV) build app_craftplace
	docker compose -v -f $(DC_DEV) up  app_craftplace -d

.PHONY: down_app
down_app:
	docker compose -f $(DC_DEV) down -v app_craftplace

.PHONY: run_rep
run_rep:
	docker compose -v -f $(DC_DEV) up --build app_craftplace_replica1 app_craftplace_replica2 -d

.PHONY: down_rep
down_rep:
	docker compose -f $(DC_DEV) down -v app_craftplace_replica1 app_craftplace_replica2

.PHONY: run_mirror
run_mirror:
	docker compose -v -f $(DC_DEV) up --build app_craftplace_mirror -d

.PHONY: down_mirror
down_mirror:
	docker compose -f $(DC_DEV) down -v app_craftplace_mirror

.PHONY: run_db
run_db:
	docker compose -v -f $(DC_DEV) --env-file $(DB_ENV) up --build postgres_craftplace -d
	./scripts/setup-replication-manually.sh 
	docker compose -v -f $(DC_DEV) --env-file $(DB_ENV) up --build pg_migrator_craftplace -d
	docker compose -v -f $(DC_DEV) --env-file $(DB_ENV) up --build postgres_craftplace_slave -d

.PHONY: down_db
down_db:
	docker compose -f $(DC_DEV) --env-file $(DB_ENV) down -v postgres_craftplace postgres_craftplace_slave pg_migrator_craftplace


# .PHONY: run_slave
# run_slave:
# 	docker compose -v -f $(DC_DEV) --env-file $(DB_ENV) up --build postgres_craftplace_slave -d

# .PHONY: down_slave
# down_slave:
# 	docker compose -v -f $(DC_DEV) --env-file $(DB_ENV) down -v postgres_craftplace_slave






.PHONY: restart_ng
restart_ng:
	docker compose -f $(DC_DEV) --env-file $(DB_ENV) restart nginx

.PHONY: run_ng
run_ng:
	docker compose -v -f $(DC_DEV) --env-file $(DB_ENV) up --build nginx -d

.PHONY: down_ng
down_ng:
	docker compose -f $(DC_DEV) --env-file $(DB_ENV) down -v nginx


.PHONY: run_all
run_all:
	docker compose -v -f $(DC_DEV) --env-file $(DB_ENV) up --build -d

.PHONY: down_all
down_all:
	docker compose -f $(DC_DEV) --env-file $(DB_ENV) down -v


.PHONY: run_pgadmin
run_pgadmin:
	docker compose -v -f $(DC_DEV) --env-file $(DB_ENV) up --build pgadmin -d

.PHONY: down_pgadmin
down_pgadmin:
	docker compose -f $(DC_DEV) --env-file $(DB_ENV) down pgadmin

.PHONY: swagger
swagger:
	swag init -g ./cmd/dev/main.go --output ./docs
	swag init -g ./cmd/mirror/main.go --output ./docs_mirror

.PHONY: docs_clear
docs_clear:
	sudo rm -rf ./docs/*
	sudo rm -rf ./tmp/*
	make swagger

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
# 	docker-compose -f ./deployment/docker-compose.ci.yml down
# Удалите старые образы
	docker rmi deployment-test-runner
# Очистите builder кеш
	docker builder prune -f
# Удалите все старые версии
	docker rmi deployment-test-runner:latest
# Полная очистка
	docker system prune -a -f



