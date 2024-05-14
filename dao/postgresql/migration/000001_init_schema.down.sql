-- 类型
drop type if exists Gender cascade;
drop type if exists RelationType cascade;
drop type if exists GroupType cascade;
drop type if exists FriendType cascade;
drop type if exists ApplicationStatus cascade;
drop type if exists FileType cascade;
drop type if exists MsgNotifyType cascade;
-- 删除表
drop type if exists users cascade;
drop table if exists accounts cascade;
drop table if exists relations cascade;
drop table if exists settings cascade;
drop table if exists applations cascade;
drop table if exists files cascade;
drop table if exists messages cascade;
drop table if exists group_notify cascade;
-- 方法
drop function if exists pin_timestamp() cascade;
drop function if exists cs_timestamp() cascade;
drop function if exists show_timestamp() cascade;
-- 触发器
drop trigger if exists message_mag_content_tsv on messages cascade;
drop trigger if exists group_notify_msg_content_tsv on group_notify cascade;
drop trigger if exists pin_timestamp_relations_settings_trigger on settings cascade;
drop trigger if exists pin_timestamp_messages_trigger on messages cascade;
drop trigger if exists application_update_at_trigger on applications cascade;
drop trigger if exists show_timestamp_trigger on settings cascade;
-- 语言
drop text search configuration if exists chinese;