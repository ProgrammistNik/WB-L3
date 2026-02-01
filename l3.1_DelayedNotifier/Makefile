run:
	docker-compose up -d

test:
	go test -cover ./...
	
race-test:
	go test -race ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: vet lint race-test test