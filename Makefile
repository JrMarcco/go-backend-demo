
DB_URL=mysql://root:u2E3WWtgam@tcp(127.0.0.1:31964)/simple_bank

migrateup:
	migrate --path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate --path db/migration -database "$(DB_URL)" -verbose down

test:
	go test -v -cover ./...

.PHONY: migrateup migratedown