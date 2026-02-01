KAFKA_TOPIC = images
KAFKA_BROKER = localhost:9092
UPLOADS_DIR = ./uploads

.PHONY: run create-topic list-topics start init-folders

init-folders:
	mkdir -p $(UPLOADS_DIR)/originals
	mkdir -p $(UPLOADS_DIR)/processed
	mkdir -p $(UPLOADS_DIR)/thumbnails

run: init-folders
	docker-compose up -d

create-topic:
	docker exec kafka kafka-topics --create \
		--topic $(KAFKA_TOPIC) \
		--bootstrap-server $(KAFKA_BROKER) \
		--partitions 1 \
		--replication-factor 1 || true

list-topics:
	docker exec kafka kafka-topics --list --bootstrap-server $(KAFKA_BROKER)

start: run create-topic list-topics
