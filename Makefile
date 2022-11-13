DB_URL=mysql://root:@tcp(127.0.0.1:3306)/simple_bank

mysql:
	docker run --name mysql -p 3306:3306 -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -d mysql:latest

createdb:
	docker exec -it mysql mysql --host 127.0.0.1 --port 3306 -uroot -e "create database if not exists simple_bank collate utf8mb4_general_ci"

migrateup:
	migrate --path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate --path db/migration -database "$(DB_URL)" -verbose down

test:
	go test -v -cover ./...

run:
	go run main.go

.PHONY: mysql createdb migrateup migratedown test run