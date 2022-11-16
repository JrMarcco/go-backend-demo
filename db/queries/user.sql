-- name: CreateUser :execresult
insert into users (username, email, hashed_passwd) values (?, ?, ?);

-- name: GetUser :one
select * from users
where id = ? limit 1;