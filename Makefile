build:
	docker build -t flex/rabbit-amazon-forwarder -f Dockerfile .

push: test build
	docker push flex/rabbit-amazon-forwarder

test:
	docker-compose run --rm tests

dev:
	go build
