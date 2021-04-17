SHELL = /bin/bash

include .env

build-bin:
	CGO_ENABLED=0 go build -a -o social-network cmd/social-network.go

build-docker:
	docker build -t social-network . -f build/Dockerfile

run-front-local:
	cd frontend && API_HOST=http://localhost:8081 npm start

social-network: build-bin
run-back-local: social-network
	go run cmd/social-network.go

run:
	docker run --env-file .env -p 8080:8080 social-network

migrate:
	cd migrations && \
		~/go/bin/goose mysql "${MYSQL_LOGIN}:${MYSQL_PASSWORD}@${MYSQL_HOST}/${MYSQL_DATABASE}" up

clean:
	rm social-network
