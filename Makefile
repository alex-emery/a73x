VERSION=0.0.1
DOCKER_IMAGE=alexemery/home:${VERSION}

.PHONY: fly 
fly: 
	cd proxy && fly deploy

.PHONY: serve 
serve: 
	go run ./cmd/home

.PHONY: image
image:
	docker build . -t ${DOCKER_IMAGE}

.PHONY: public
public: image
	docker push ${DOCKER_IMAGE}
