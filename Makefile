.PHONY: build
build:
	docker build -t poncheska/sa-storage-v2 -f builds/Dockerfile .
	docker push poncheska/sa-storage-v2

.PHONY: run
run:
	go run ./main.go
