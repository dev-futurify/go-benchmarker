DOCKER_IMAGE = "go-benchmarker"
DOCKER_HUB_IMAGE = "ricogustavo/go-benchmarker"


docker.build:
	docker build -t $(DOCKER_IMAGE) . --platform linux/amd64

docker.run: docker.build
	docker run -it --rm $(DOCKER_IMAGE)

docker.push: docker.build
	docker tag $(DOCKER_IMAGE) $(DOCKER_HUB_IMAGE):latest
	docker push $(DOCKER_HUB_IMAGE):latest