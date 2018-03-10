PORT=8080
build:
	go build backend.go
test:
	docker system prune -f
	docker-compose -f Docker-compose.test.yml up --build --abort-on-container-exit
	# docker build -t backend_tests -f Dockerfile.test .
	# docker run backend_tests
run-main:
	go build main.go
	./main

run:
	sudo cp /etc/letsencrypt/live/engineerbeard.com/privkey.pem .
	sudo cp /etc/letsencrypt/live/engineerbeard.com/fullchain.pem .
	docker system prune -f
	docker-compose -f Docker-compose.run.yml up --build --abort-on-container-exit

dev:
	docker system prune -f
	docker-compose -f Docker-compose.run.dev.yml up --build --abort-on-container-exit

run-d:
	docker system prune -f
	docker-compose -f Docker-compose.run.yml up --build -d

stop:
	docker stop barrenschatapi_barrenschat-api-1_1
	docker stop barrenschatapi_barrenschat-api-2_1
	docker stop barrenschatapi_load-balancer_1
	docker stop barrenschatapi_redis_1

stop-all:
	docker stop $(docker ps -aq)

certs:
	sudo ~/certbot-auto certonly
	sudo cp /etc/letsencrypt/live/engineerbeard.com/privkey.pem .
	sudo cp /etc/letsencrypt/live/engineerbeard.com/fullchain.pem .