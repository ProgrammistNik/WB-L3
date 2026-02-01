run:
	docker-compose up -d

race-test:
	go test -race ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: vet lint race-test