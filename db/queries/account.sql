-- name: CreateAccount :execresult
insert into account (account_owner, balance, currency) values (?, ?, ?);

-- name: GetAccount :one
select * from account
where id = ? limit 1;

-- name: ListAccount :many
select * from account
order by id limit ?, ?;

-- name: AddBalance :exec
update account set balance = balance + sqlc.arg(amount) where id = ?;

-- name: DeleteAccount :exec
delete from account where id = ? limit 1;