PORT=8080
build:
	go build backend.go
test:
	docker-compose -f Docker-compose.test.yml up --build --abort-on-container-exit
	# docker build -t backend_tests -f Dockerfile.test .
	# docker run backend_tests
docker-run:
	docker build -t backend -f Dockerfile.run .
	docker run -it backend -p "8080:8080"