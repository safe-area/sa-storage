.PHONY: build
build:
	docker build -t poncheska/sa-storage -f builds/Dockerfile .
	docker push poncheska/sa-storage

.PHONY: run
run:
	go run ./main.go
