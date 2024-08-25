.PHONY: fly 

fly: 
	cd proxy && fly deploy

.PHONY: content 
content:
	go run ./cmd/generate

.PHONY: serve
serve: content
	go run ./cmd/home
