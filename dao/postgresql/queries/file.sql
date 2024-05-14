-- name: CreateFile :one
insert into files (file_name, file_type, file_size, key, url, relation_id, account_id)
values ($1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: DeleteFileByID :exec
delete
from files
where id = $1;

-- name: GetFileKeyByID :one
select key
from files
where id = $1;

-- name: GetFileByRelationID :many
select *
from files
where relation_id = $1;

-- name: GetFileDetailsByID :one
select *
from files
where id = $1;

-- name: GetGroupAvatar :one
select *
from files
where relation_id = $1
    and account_id is null;

-- name: UpdateGroupAvatar :exec
update files
set url = $1
where relation_id = $2 and file_name = 'groupAvatar';

-- name: GetFileByRelationIDIsNULL :many
select id, key
from files
where relation_id is null and file_name != 'AccountAvatar';
