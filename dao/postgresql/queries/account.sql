-- name: CreateAccount :exec
insert into accounts (id, user_id, name, avatar, gender, signature)
values ($1, $2, $3, $4, $5, $6);

-- name: DeleteAccount :exec
delete
from accounts
where id = $1;

-- name: DeleteAccountsByUserID :many
delete
from accounts
where user_id = $1
returning id;

-- name: UpdateAccount :exec
update accounts
set name = $1,
    gender = $2,
    signature = $3
where id = $4;

-- name: UpdateAccountAvatar :exec
update accounts
set avatar = $1
where id = $2;

-- name: GetAccountByID :one
select a.*, r.id as relation_id
from (select * from accounts where accounts.id = @target_id) a
    left join relations r on
        r.relation_type = 'friend' and
        ((r.friend_type).account1_id = a.id and (r.friend_type).account2_id = @self_id::bigint) or
        (r.friend_type).account1_id = @self_id::bigint and (r.friend_type).account2_id = a.id
limit 1;

-- name: GetAccountsByUserID :many
select id, name, avatar, gender
from accounts
where user_id = $1;

-- name: ExistsAccountByID :one
select exists(
           select 1
           from accounts
           where id = $1
);

-- name: ExistsAccountByNameAndUserID :one
select exists(
    select 1
    from accounts
    where user_id = $1
    and name = $2
);

-- name: CountAccountsByUserID :one
select count(id)::int
from accounts
where user_id = $1;

-- name: GetAccountsByName :many
select a.*, r.id as relation_id, count(*) over () as total
from (select id, name, avatar, gender from accounts where name like ('%' || @name::varchar || '%')) as a
    left join relations r on (r.relation_type = 'friend' and
                             (((r.friend_type).account1_id = a.id and (r.friend_type).account2_id = @account_id::bigint) or
                              ((r.friend_type).account1_id = @account_id::bigint and (r.friend_type).account2_id = a.id)))
limit $1 offset $2;
