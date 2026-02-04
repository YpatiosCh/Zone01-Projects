## run the program
run:
	go run cmd/web/main.go

# For Docker
docker-build:
	docker build -t forum-app .
# For Docker-compose
docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down
# For Podman-compose 
podman-run:
	podman-compose up -d

podman-stop:
	podman-compose down

# Build script (simplifies the build process)
build-docker:
	chmod +x build.sh
	./build.sh

# Cleanup unused Docker objects
cleanup:
	chmod +x cleanup.sh
	./cleanup.sh
