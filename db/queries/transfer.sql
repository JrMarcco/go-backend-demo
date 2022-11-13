-- name: CreateTransfer :execresult
insert into transfer (from_id, to_id, amount) values (?, ?, ?);
