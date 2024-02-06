build:
	docker build -t davidmwhynot/rabbit-amazon-forwarder -f Dockerfile .

push: test build
	docker push davidmwhynot/rabbit-amazon-forwarder

test:
	docker-compose run --rm --build tests

dev:
	go build
