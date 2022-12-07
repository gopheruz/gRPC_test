-include .env
.SILENT:
DB_URL="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable"	
CURRENT_DIR=$(shell pwd)
proto-gen:
	rm -rf genproto
	./gen-proto.sh ${CURRENT_DIR}


migrateup:
	migrate -path migrations -database $(DB_URL) -verbose up


migratedown:
	migrate -path migrations -database $(DB_URL) -verbose down
