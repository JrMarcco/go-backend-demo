
DB_URL=mysql://root:u2E3WWtgam@tcp(192.168.3.50:31964)/simple_bank

migrateup:
	migrate --path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate --path db/migration -database "$(DB_URL)" -verbose down

test:
	go test -v -cover ./...

.PHONY: migrateup migratedown