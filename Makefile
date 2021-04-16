build-bin:
	CGO_ENABLED=0 go build -a -o social-network cmd/social-network.go

build-docker:
	docker build -t social-network . -f build/Dockerfile

run:
	docker run -e PORT=8080 -e HOST=0.0.0.0 -p 8080:8080 social-network

clean:
	rm social-network
