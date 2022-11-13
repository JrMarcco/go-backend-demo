migrateup:
	migrate --path db/migration -database "mysql://root:u2E3WWtgam@tcp(192.168.3.50:31964)/simple_bank" -verbose up

migratedown:
	migrate --path db/migration -database "mysql://root:u2E3WWtgam@tcp(192.168.3.50:31964)/simple_bank?" -verbose down

.PHONY: migrateup migratedown