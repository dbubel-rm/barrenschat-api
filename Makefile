PORT=8080
build:
	go build backend.go
test:
	docker-compose -f Docker-compose.test.yml up --build --abort-on-container-exit
	# docker build -t backend_tests -f Dockerfile.test .
	# docker run backend_tests
run:
	go build main.go
	./main

run-docker:
	docker-compose -f Docker-compose.run.yml up --build --abort-on-container-exit
	