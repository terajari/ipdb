include .env

run/api:
	go run ./cmd/api/

docker/up:
	docker-compose up -d

docker/down:
	docker-compose down

migrate/create:
	migrate create -ext sql -dir ./migrations -seq ${name}

migrate/up:
	migrate -path ./migrations -database "${DB_DSN}" -verbose up

migrate/down:
	migrate -path ./migrations -database "${DB_DSN}" -verbose down

migrate/fix:
	migrate -path ./migrations -database "${DB_DSN}" -verbose force ${version}