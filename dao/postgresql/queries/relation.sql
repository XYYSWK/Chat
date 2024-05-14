-- name: CreateGroupRelation :one
insert into relations (relation_type, group_type)
values ('group', ROW (@name::varchar(50), @description::varchar(255), @avatar::varchar(255)))
returning id;

-- name: CreateFriendRelation :one
insert into relations (relation_type, friend_type)
values ('friend', ROW (@account1_id::bigint, @account2_id::bigint))
returning id;

-- name: DeleteRelation :exec
delete
from relations
where id = @id;

-- name: DeleteFriendRelationsByAccountID :many
delete
from relations
where relation_type = 'friend'
 and ((friend_type).account1_id = @account1_id::bigint or (friend_type).account2_id = @account1_id::bigint)
returning id;

-- name: UpdateGroupRelation :exec
update relations
set group_type = (ROW (@name::varchar(50), @description::varchar(255), @avatar::varchar(255)))
where relation_type = 'group'
 and id = @id;

-- name: GetGroupRelationByID :one
select id,
        relation_type,
        (group_type).name::varchar as name,
        (group_type).description::varchar as description,
        (group_type).avatar::varchar as avatar,
        create_at
from relations
where relation_type = 'group'
 and id = @id;

-- name: ExistsFriendRelation :one
select exists(select 1
              from relations
              where relation_type = 'friend'
               and (friend_type).account1_id = @account1_id::bigint
               and (friend_type).account2_id = @account2_id::bigint);

-- name: GetFriendRelationByID :one
select (friend_type).account1_id as account1_id,
       (friend_type).account2_id as account2_id,
       create_at
from relations
where relation_type = 'friend'
    and id = $1;

-- name: GetAllGroupRelation :many
select id
from relations
where relation_type = 'group'
    and friend_type is null;

-- name: GetAllRelationOnRelation :many
select *
from relations;

-- name: GetAllRelationIDs :many
select id
from relations;

-- name: GetRelationIDByAccountID :one
select id
from relations
where (friend_type).account1_id = @account_id::bigint
  and (friend_type).account2_id = @account_id::bigint;