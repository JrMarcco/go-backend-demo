-- name: CreateEntry :execresult
insert into entry (account_id, amount) values (?, ?);
