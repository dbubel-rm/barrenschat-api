PORT=8080
build:
	go build backend.go
test:
	docker-compose -f Docker-compose.test.yml up --build --abort-on-container-exit
	# docker build -t backend_tests -f Dockerfile.test .
	# docker run backend_tests
run-main:
	go build main.go
	./main

app-run:
	docker system prune -f
	docker-compose -f Docker-compose.run.yml up --build --abort-on-container-exit

app-stop:
	docker stop barrenschatapi_barrenschat-api-1_1
	docker stop barrenschatapi_barrenschat-api-2_1
	docker stop barrenschatapi_load-balancer_1
	docker stop barrenschatapi_redis_1

stop-all:
	docker stop $(docker ps -aq)
