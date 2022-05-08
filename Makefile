SHELL = /bin/bash

include .env

build-bin: social-network
build-migrator: migrator

build-docker:
	docker build -t social-network . -f build/Dockerfile

run-front-local:
	cd frontend && API_HOST=http://localhost:8081 npm start

social-network:
	CGO_ENABLED=0 go build -a -o social-network cmd/social-network/social-network.go
run-back-local: social-network
	./social-network

run:
	docker run --env-file .env -p ${PORT}:${PORT} social-network

migrator:
	CGO_ENABLED=0 go build -a -o migrator cmd/migrator/migrator.go
migrate: migrator
	./migrator -dir migrations/ "${MYSQL_USER}:${MYSQL_PASSWORD}@${MYSQL_HOST}/${MYSQL_DATABASE}?parseTime=true" up
migrate-down: migrator
	./migrator -dir migrations/ "${MYSQL_USER}:${MYSQL_PASSWORD}@${MYSQL_HOST}/${MYSQL_DATABASE}?parseTime=true" down

clean-bin:
	rm social-network || true
clean-migrator:
	rm migrator || true

clean: clean-bin clean-migrator
