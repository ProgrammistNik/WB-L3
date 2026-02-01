run:
	docker-compose up -d
	
vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: vet lint

clean:
	docker-compose down --rmi all --volumes --remove-orphans