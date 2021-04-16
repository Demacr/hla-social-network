SHELL = /bin/bash

build-bin:
	CGO_ENABLED=0 go build -a -o social-network cmd/social-network.go

build-docker:
	docker build -t social-network . -f build/Dockerfile

run-front-local:
	cd frontend && API_HOST=http://localhost:8081 npm start

run-back-local:
	set -a && source .env && go run cmd/social-network.go

run:
	docker run --env-file .env -p 8080:8080 social-network

clean:
	rm social-network
