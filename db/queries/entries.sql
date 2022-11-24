-- name: CreateEntries :execresult
insert into entries (account_id, amount) values (?, ?);
