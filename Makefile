name = wp-atrd-task
service_host = 127.0.0.1
service_port = 8080
doc_port = 8081

.DEFAULT_GOAL := $(name)

all: $(name) doc test clean

clean: clean-doc clean-$(name)

$(name): $(name)-image network $(name)-database $(name)-container

clean-$(name): stop-$(name)-container clean-$(name)-container stop-$(name)-database clean-$(name)-database clean-network clean-$(name)-image

clean-$(name)-image: clean-$(name)-container
	docker image inspect $(name) &> /dev/null && docker rmi $(name):latest || true

$(name)-image: clean-$(name)-image
	docker build -t $(name) .

$(name)-container: clean-$(name)-container
	docker run -d --name $(name) \
		--network=$(name) \
		-p 127.0.0.1:$(service_port):8080 \
		$(name)

clean-$(name)-container: stop-$(name)-container
	docker container inspect $(name) &> /dev/null && docker container rm $(name) || true

stop-$(name)-container:
	docker container inspect $(name) &> /dev/null && docker container stop $(name) || true

network: clean-network
	docker network create $(name)

clean-network: clean-$(name)-database clean-doc
	docker network inspect $(name) &> /dev/null && docker network rm $(name) || true

network-cli: network
	docker run -it --net $(name) nicolaka/netshoot

$(name)-database: clean-$(name)-database
	docker run --name $(name)-database \
		--network=$(name) \
		-d redis
		# -v $(shell pwd)/redis:/data \

stop-$(name)-database:
	docker container inspect $(name)-database &> /dev/null && docker container stop $(name)-database || true

clean-$(name)-database: stop-$(name)-database
	docker container inspect $(name)-database &> /dev/null && docker container rm $(name)-database || true

$(name)-database-cli:
	docker exec -it $(name)-database redis-cli



test: $(name)-test-image $(name)-test-container

clean-test: stop-$(name)-test-container clean-$(name)-test-container clean-$(name)-test-image

$(name)-test-image: clean-$(name)-test-image
	docker build -t $(name)-test -f ./Dockerfile.test .

$(name)-test-container: clean-$(name)-test-container
	docker run --name $(name)-test \
		--network=$(name) \
		$(name)-test

test-cli: clean-$(name)-test-container
	docker run -it --name $(name)-test \
		--network=$(name) \
		$(name)-test \
		/bin/bash

clean-$(name)-test-image: clean-$(name)-test-container 
	docker image inspect $(name)-test &> /dev/null && docker rmi $(name)-test:latest || true

clean-$(name)-test-container: stop-$(name)-test-container
	docker container inspect $(name)-test &> /dev/null && docker container rm $(name)-test || true

stop-$(name)-test-container:
	docker container inspect $(name)-test &> /dev/null && docker container stop $(name)-test || true



godog: $(name)-godog-image $(name)-godog-container

clean-godog: stop-$(name)-godog-container clean-$(name)-godog-container clean-$(name)-godog-image

$(name)-godog-image: clean-$(name)-godog-image
	docker build -t $(name)-godog -f ./Dockerfile.godog .

$(name)-godog-container: clean-$(name)-godog-container
	docker run --name $(name)-godog \
		--network=$(name) \
		$(name)-godog

clean-$(name)-godog-image: clean-$(name)-godog-container
	docker image inspect $(name)-godog &> /dev/null && docker rmi $(name)-godog:latest || true

clean-$(name)-godog-container: stop-$(name)-godog-container
	docker container inspect $(name)-godog &> /dev/null && docker container rm $(name)-godog || true

stop-$(name)-godog-container:
	docker container inspect $(name)-godog &> /dev/null && docker container stop $(name)-godog || true

godog-cli:
	docker run -rm -it --name $(name)-godog \
		--network=$(name) \
		$(name)-godog \
		/bin/bash


doc: start-$(name)-doc-container

clean-doc: clean-$(name)-doc-container

start-$(name)-doc-container: clean-$(name)-doc-container
	docker run -d --name $(name)-doc \
		--network=$(name) \
		-p $(doc_port):8080 \
		-e URL=http://$(service_host):$(service_port)/swagger.yml \
		swaggerapi/swagger-ui

clean-$(name)-doc-container: stop-$(name)-doc-container
	docker container inspect $(name)-doc &> /dev/null && docker container rm $(name)-doc || true

stop-$(name)-doc-container:
	docker container inspect $(name)-doc &> /dev/null && docker container stop $(name)-doc || true
