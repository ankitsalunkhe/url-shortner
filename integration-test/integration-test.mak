.PHONY: integration-test
integration-test: integration-test-up integration-test-down

.PHONY: integration-test-up
integration-test:
	docker compose -f ./integration-test/docker-compose.yaml up

.PHONY: integration-test-down
integration-test:
	docker compose -f ./integration-test/docker-compose.yaml down