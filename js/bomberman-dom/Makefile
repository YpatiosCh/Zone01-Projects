run:
	go run ./server/cmd/.

docker-run:
	docker build -t bomberman . && docker run --rm -p 8080:8080 bomberman

docker-build:
	docker build -t bomberman .

upload:
	docker build -t bomberman .
	docker tag bomberman aleksisgioldaseas/bomberman:latest
	docker push aleksisgioldaseas/bomberman:latest