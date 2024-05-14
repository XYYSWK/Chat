-- name: CreateUser :one
insert into users (email, password)
values ($1, $2)
returning *;

-- name: DeleteUser :exec
delete
from users
where id = $1;

-- name: GetUserByEmail :one
select *
from users
where email = $1
limit 1;

-- name: GetUserByID :one
select *
from users
where id = $1
limit 1;

-- name: UpdateUser :exec
update users
set email = $1,
    password = $2
where id = $3;

-- name: ExistEmail :one
select exists(select 1 from users where email = $1);

-- name: GetAllEmail :many
select email
from users;

-- name: ExistsUserByID :one
select exists(select 1 from users where id = $1);

-- name: GetAcountIDsByUserID :many
select id
from accounts
where user_id = $1;
