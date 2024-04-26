devdb:
	docker-compose -f docker-compose.dev.yml up -d
devdbdown:
	docker-compose -f docker-compose.dev.yml down
migratecreate:
	migrate create -ext sql -dir db/migrations -seq $(name)
migrateup:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/temp_db?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/temp_db?sslmode=disable" -verbose down