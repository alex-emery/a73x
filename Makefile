VERSION=0.0.1
DOCKER_IMAGE=a73xsh/home:${VERSION}


.PHONY: build 
build: 
	go run ./cmd/build

.PHONY: serve 
serve: 
	go run ./cmd/serve

.PHONY: image
image:
	docker build . -t ${DOCKER_IMAGE}

.PHONY: public
public: image
	docker push ${DOCKER_IMAGE}
