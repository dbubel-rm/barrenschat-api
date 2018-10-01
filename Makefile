PORT=8080
build:
	go build -x -a main.go

run:
	sudo cp /etc/letsencrypt/live/engineerbeard.com/privkey.pem .
	sudo cp /etc/letsencrypt/live/engineerbeard.com/fullchain.pem .
	#docker rmi $(docker images -f "dangling=true" -q)
	docker system prune -f
	docker-compose -f Docker-compose.run.yml up --build --abort-on-container-exit

dev:
	#docker system prune --volumes -f
	docker-compose -f Docker-compose.run.dev.yml up --build --abort-on-container-exit

run-d:
	sudo cp /etc/letsencrypt/live/engineerbeard.com/privkey.pem .
	sudo cp /etc/letsencrypt/live/engineerbeard.com/fullchain.pem .
	#docker rmi $(docker images -f "dangling=true" -q)
	docker system prune --volumes -f
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
	
tests:
	docker system prune -f
	docker-compose -f Docker-compose.test.yml up --build --abort-on-container-exit
