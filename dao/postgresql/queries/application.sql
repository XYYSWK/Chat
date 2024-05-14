-- name: CreateApplication :exec
insert into applications (account1_id, account2_id, apply_msg, refuse_msg)
VALUES ($1, $2, $3, '');

-- name: ExistsApplicationByIDWithLock :one
select exists(
                select 1
                from applications
                where (account1_id = $1 and account2_id = $2)
                or (account1_id = $2 and account2_id = $1)
                for update );

-- name: DeleteApplication :exec
delete
from applications
where account1_id = $1 and account2_id = $2;

-- name: GetApplicationByID :one
select *
from applications
where account1_id = $1 and account2_id = $2
limit 1;

-- name: UpdateApplication :exec
update applications
set status = $2,
    refuse_msg = $1
where account1_id = $3
  and account2_id = $4;

-- name: GetApplications :many
select app.*,
       a1.name as account1_name,
       a1.avatar as account1_avatar,
       a2.name as account2_name,
       a2.avatar as account2_avatar
from accounts a1,
     accounts a2,
     (select *, count(*) over () as total
      from applications
      where account1_id = @account_id
      or account2_id = @account_id
      order by create_at desc
      limit $1 offset $2) as app
where a1.id = app.account1_id
and a2.id = app.account2_id;