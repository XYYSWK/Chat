-- name: CreateGroupNotify :one
insert into group_notify
(relation_id, msg_content, msg_expand, account_id, create_at, read_ids)
values ($1, $2, $3, $4, $5, $6)
returning id, relation_id, msg_content, msg_expand, account_id, create_at, read_ids;

-- name: UpdateGroupNotify :one
update group_notify
set relation_id = $1,
    msg_content = $2,
    msg_expand = $3,
    account_id = $4,
    create_at = $5,
    read_ids = $6
where id = $7
returning id, relation_id, msg_content, msg_expand, account_id, create_at, read_ids;

-- name: GetGroupNotifyByID :many
select id, relation_id, msg_content, msg_expand, account_id, create_at, read_ids
from group_notify
where relation_id = $1;

-- name: DeleteGroupNotify :exec
delete
from group_notify
where id = $1;