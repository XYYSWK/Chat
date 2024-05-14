-- name: CreateSetting :exec
insert into settings (account_id, relation_id, nick_name, is_leader, is_self)
values ($1, $2, '', $3, $4);

-- name: DeleteSetting :exec
delete
from settings
where account_id = $1
 and relation_id = $2;

-- name: DeleteSettingsByAccountID :many
delete
from settings
where account_id = $1
returning relation_id;

-- name: ExistsGroupLeaderByAccountIDWithLock :one
select exists(select 1 from settings where account_id = $1 and is_leader = true) for update;

-- name: UpdateSettingNickName :exec
update settings
set nick_name = $1
where account_id = $2
 and relation_id = $3;

-- name: UpdateSettingDisturb :exec
update settings
set is_not_disturb = $1
where account_id = $2
    and relation_id = $3;

-- name: UpdateSettingPin :exec
update settings
set is_pin = $1
where account_id = $2
    and relation_id = $3;

-- name: UpdateSettingLeader :exec
update settings
set is_leader = $1
where account_id = $2
    and relation_id = $3;

-- name: UpdateSettingShow :exec
update settings
set is_show = $1
where account_id = $2
    and relation_id = $3;

-- name: GetSettingByID :one
select *
from settings
where account_id = $1
    and relation_id = $2;

-- name: GetFriendPinSettingsOrderByPinTime :many
select s.*,
       a.id as account_id,
       a.name as account_name,
       a.avatar as account_avatar
from (select settings.relation_id, settings.nick_name, settings.pin_time
      from settings,
           relations
      where settings.account_id = $1
        and settings.is_pin = true
        and settings.relation_id = relations.id
        and relations.relation_type = 'friend') as s,
    accounts a
where a.id = (select account_id from settings where relation_id = s.relation_id and (account_id != $1 or is_self = true))
order by s.pin_time;

-- name: GetGroupPinSettingsOrderByPinTime :many
select s.relation_id,
       s.nick_name,
       s.pin_time,
       r.id,
       r.group_type
from (select settings.relation_id, settings.nick_name, settings.pin_time
      from settings,
           relations
      where settings.account_id = $1
        and settings.relation_id = relations.id
        and settings.is_pin = true
        and relation_type = 'group') as s,
    relations r
where r.id = (select relation_id from settings where relation_id = s.relation_id and account_id = $1)
order by s.pin_time;

-- name: GetFriendShowSettingsOrderByShowTime :many
select s.*,
       a.id as account_id,
       a.name as account_name,
       a.avatar as account_avatar
from (select relation_id,
             nick_name,
             is_not_disturb,
             is_pin,
             pin_time,
             is_show,
             last_show,
             is_self
      from settings,
           relations
      where settings.account_id = $1
        and settings.is_show = true
        and settings.relation_id = relations.id
        and relations.relation_type = 'friend') as s,
    accounts a
where a.id = (select account_id from settings where relation_id = s.relation_id and (account_id != $1 or is_self = true))
order by s.last_show desc;

-- name: GetGroupShowSettingsOrderByShowTime :many
select s.*,
       r.id,
       r.group_type
from (select relation_id,
             nick_name,
             is_not_disturb,
             is_pin,
             pin_time,
             is_show,
             last_show,
             is_self
             from settings,
                  relations
             where settings.account_id = $1
                and settings.relation_id = relations.id
                and settings.is_show = true
                and relations.relation_type = 'group') as s,
    relations r
where r.id = (select relation_id from settings where relation_id = s.relation_id and account_id = $1)
order by s.last_show desc;

-- name: GetFriendSettingsOrderByName :many
select s.*,
       a.id as account_id,
       a.name as account_name,
       a.avatar as account_avatar
from (select relation_id,
             nick_name,
             is_not_disturb,
             is_pin,
             pin_time,
             is_show,
             last_show,
             is_self
             from settings,
                  relations
             where settings.account_id = $1
                and settings.relation_id = relations.id
                and relations.relation_type = 'friend') as s,
    accounts a
where a.id = (select account_id from settings where relation_id = s.relation_id and (account_id != $1 or s.is_self = true))
order by a.name;

-- name: ExistsFriendSetting :one
select exists(select 1
              from settings s,
                   relations r
              where r.relation_type = 'friend'
                and (((r.friend_type).account1_id = @account1_id::bigint and
                     (r.friend_type).account2_id = @account2_id::bigint) or
                     ((r.friend_type).account1_id = @account2_id::bigint and
                      (r.friend_type).account2_id = @account1_id::bigint) )
                and s.account_id = @account1_id
              );

-- name: GetFriendSettingsByName :many
select s.*,
       a.id as account_id,
       a.name as account_name,
       a.avatar as account_avatar,
       count(*) over () as total
from (select relation_id,
             nick_name,
             is_not_disturb,
             is_pin,
             pin_time,
             is_show,
             last_show,
             is_self
             from settings,
                  relations
             where settings.account_id = $1
                and settings.relation_id = relations.id
                and relations.relation_type = 'friend') as s,
    accounts a
where a.id = (select account_id from settings where relation_id = s.relation_id and (account_id != $1 or s.is_self = true))
    and ((a.name like ('%' || @name::varchar || '%')) or (nick_name like ('%' || @name::varchar || '%')))
order by a.name
limit $2 offset $3;

-- name: GetGroupSettingsByName :many
select s.*,
       r.id as relation_id,
       (r.group_type).name as group_name,
       (r.group_type).avatar as group_avatar,
       (r.group_type).description as description,
       count(*) over () as total
from (select relation_id,
             nick_name,
             is_not_disturb,
             is_pin,
             pin_time,
             is_show,
             last_show,
             is_self
          from settings,
               relations
          where settings.account_id = $1
            and settings.relation_id = relations.id
            and relations.relation_type = 'group') as s,
    relations r
where r.id = (select s.relation_id from settings where relation_id = s.relation_id and (settings.account_id = $1))
    and (((r.group_type).name like ('%' || @name::varchar || '%')))
order by (r.group_type).name
limit $2 offset $3;

-- name: TransferIsLeaderTrue :exec
update settings
set is_leader = true
where relation_id = $1
    and account_id = $2;

-- name: TransferIsLeaderFalse :exec
update settings
set is_leader = false
where relation_id = $1
    and account_id = $2;

-- name: DeleteGroup :exec
delete
from settings
where relation_id = $1;

-- name: ExistsSetting :one
select exists(select 1 from settings where account_id = $1 and relation_id = $2);

-- name: ExistsIsLeader :one
select exists(select 1 from settings where relation_id = $1 and account_id = $2 and is_leader is true);

-- name: GetGroupMembers :many
select account_id
from settings
where relation_id = $1;

-- name: GetAccountIDsByRelationID :many
select distinct account_id
from settings
where relation_id = $1;

-- name: GetGroupList :many
select s.*,
       r.id as relation_id,
       (r.group_type).name as group_name,
       (r.group_type).description as discription,
       (r.group_type).avatar as group_avatar,
       count(*) over () as total
from (select relation_id,
             nick_name,
             is_not_disturb,
             is_pin,
             pin_time,
             is_show,
             last_show,
             is_self
      from settings,
           relations
      where settings.account_id = $1
        and settings.relation_id = relations.id
        and relations.relation_type = 'group') as s,
    relations r
where r.id = (select s.relation_id from settings where relation_id = s.relation_id and (settings.account_id = $1))
order by s.last_show;

-- name: CreateManySetting :copyfrom
insert into settings (account_id, relation_id, nick_name) values ($1, $2, $3);

-- name: GetGroupMembersByID :many
select a.id, a.name, a.avatar, s.nick_name, s.is_leader
from accounts a
    left join settings s on a.id = s.account_id
where s.relation_id = $1
limit $2 offset $3;