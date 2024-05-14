-- 连接到目标数据库，创建zhparser解析器
-- CREATE EXTENSION zhparser;
-- 将zhparser解析器作为全文检索配置项
-- CREATE TEXT SEARCH CONFIGURATION chinese (PARSER = zhparser);
-- 普遍情况下，我们只需要按照名词(n)，动词(v)，形容词(a)，成语(i),叹词(e)和习用语(l)6种方式对句子进行划分就可以了，词典使用的是内置的simple词典，即仅做小写转换
-- ALTER TEXT SEARCH CONFIGURATION chinese ADD MAPPING FOR n,v,a,i,e,l WITH simple;

-- 创建类型

-- 性别类型
create type Gender As ENUM ('男', '女', '未知');
-- 群或好友关系的类型
create type RelationType As ENUM ('group', 'friend');
-- 群类型
create type GroupType As
(
    name varchar(50), -- 群名称
    description varchar(255), -- 群描述
    avatar varchar(255) -- 群头像
);
-- 好友类型
create type FriendType AS
(
    account1_id bigint, -- 好友 1 的账号 id
    account2_id bigint -- 好友 2 的账号 id
);
-- 申请状态
create type ApplicationStatus AS ENUM ('已申请','已同意','已拒绝');
-- 文件类型
create type FileType As ENUM ('img','file');
-- 消息通知类型
create type MsgNotifyType As ENUM ('system', 'common');

-- 创建表

-- 用户
create table if not exists users (
    id bigserial primary key, -- 用户 id（大自增整数）
    email varchar(255) not null unique, -- 邮箱000001_init_schema.down.sql
    password varchar(255) not null, -- 密码
    create_at timestamptz not null default now() -- 创建时间（带时区）
);

-- 账号
create table if not exists accounts (
    id bigint primary key, -- 账号 id
    user_id bigint not null references users (id) on delete cascade on update cascade, -- 用户 id（外键）
    name varchar(255) not null, -- 账号名
    avatar varchar(255) not null, -- 账号头像
    gender Gender not null default '未知', -- 账号性别
    signature text not null default '这个用户很懒，什么也没有留下~', -- 账号签名
    create_at timestamptz not null default now(), -- 创建时间
    constraint account_unique_name unique (user_id, name) -- 一个用户的不同账号名不能重复
);

-- 账号名和头像索引(可以加快根据账号名和头像查询的速度)
create index account_index_name_avatar on accounts(name, avatar);

-- 群组或好友
create table relations
(
    id bigserial primary key, -- id
    relation_type RelationType not null, -- 关系类型 group:群组,friend:好友
    group_type GroupType, -- 群组信息，只有群组才有这个字段，否则为 null
    friend_type FriendType, -- 好友信息，只有好友才有这个字段，否则为 null
    create_at timestamptz default now(), -- 创建时间
    check (((group_type is null) or (friend_type is null)) and
           ((group_type is not null) or (friend_type is not null))) -- 只能存在一种信息
);

-- 账号对群组或好友关系的设置
create table settings
(
    account_id bigint not null references accounts (id) on delete cascade on update cascade, -- 账号id（外键）
    relation_id bigint not null references relations (id) on delete cascade on update cascade, -- 关系 id（外键）
    nick_name varchar(255) not null, -- 昵称，默认是账户名或群组名
    is_not_disturb boolean not null default false, -- 是否免打扰
    is_pin boolean not null default false, -- 是否置顶
    pin_time timestamptz not null default now(), -- 置顶时间
    is_show boolean not null default true, -- 是否显示
    last_show timestamptz not null default now(), -- 最后一次显示时间
    is_leader boolean not null default false, -- 是否是群主，仅对群组有效
    is_self boolean not null default false -- 是否是自己对自己的关系，仅对好友有效
);

-- 昵称索引(可以加快根据昵称查询的速度)
create index relation_setting_nickname on settings (nick_name);
-- 账户ID和关系ID的符合索引
create index setting_idx_account_id_relation_id on settings (account_id, relation_id);

-- 好友申请
create table applications
(
    account1_id bigint not null references accounts (id) on delete cascade on update cascade, -- 申请者账号 id（外键）
    account2_id bigint not null references accounts (id) on delete cascade on update cascade, -- 被申请者账号 id（外键）
    apply_msg text not null, -- 申请信息
    refuse_msg text not null, -- 拒绝信息
    status ApplicationStatus not null default '已申请', -- 申请状态
    create_at timestamptz not null default now(), -- 创建时间
    update_at timestamptz not null default now(), -- 更新时间
    constraint f_a_pk primary key (account1_id, account2_id)
);

-- 文件记录
create table files
(
    id bigserial primary key, -- 文件 id
    file_name varchar(255) not null, -- 文件名称
    file_type FileType not null, -- 文件类型
    file_size bigint not null, -- 文件大小 byte
    key varchar(255) not null, -- 文件 key 用于 obs 中删除文件
    url varchar(255) not null, -- 文件 url
    relation_id bigint references relations (id) on delete cascade on update cascade, -- 关系 id（外键）
    account_id bigint references accounts (id) on delete cascade on update cascade, -- 发送账号 id（外键）
    create_at timestamptz not null default now() -- 创建时间
);
-- 文件关系id索引
create index file_relation_id on files (relation_id);

-- 消息
create table messages
(
    id bigserial primary key, -- 消息 id
    notify_type MsgNotifyType not null, -- 消息通知类型 system:系统消息，common:普通消息
    msg_type varchar(32) not null check ( msg_type in ('text', 'file') ), -- 消息类型 text:文本消息，file:文件消息
    msg_content text not null, -- 消息内容
    msg_extend json, -- 消息扩展信息
    file_id bigint references files (id) on delete cascade on update cascade, -- 文件 id（外键），如果不是文件类型则为 null
    account_id bigint references accounts (id) on delete set null on update cascade, -- 发送账号 id（外键）
    rly_msg_id bigint references messages (id) on delete cascade on update cascade, -- 回复消息 id，没有则为 null（外键）
    relation_id bigint not null references relations (id) on delete cascade on update cascade, -- 关系 id（外键）
    create_at timestamptz not null default now(), -- 创建时间
    is_revoke boolean not null default false, -- 是否撤回
    is_top boolean not null default false, -- 是否置顶
    is_pin boolean not null default false, -- 是否pin
    pin_time timestamptz not null default now(), -- pin时间
    read_ids bigint[] not null default '{}'::bigint[], --一已读用户 id 集合
    msg_content_tsy tsvector, -- 消息分词
    check (notify_type = 'common' or (notify_type = 'system' and account_id is null)), -- 系统消息时发生账号 id 为 null
    check (msg_type = 'text' or (msg_type = 'file' and file_id is not null)) -- 文件消息时文件 id 不能为 null
);
-- 创建时间索引
create index msg_create_at on messages (create_at);
-- 分词索引(基于 GIN（Generalized Inverted Index）的全文搜索索引)
create index message_msg_content_tsv on messages using gin (to_tsvector('chinese',msg_content));

-- --
-- 触发器更新 message_msg_content_tsv
-- (在每次插入或更新 msg_content 列内容之前，通过触发器调用存储过程将文本内容转换为文本搜索向量，并存储在指定的列中，以支持全文搜索功能。)
create trigger message_mag_content_tsv
    before insert or update of msg_content
    on messages
    for each row
execute procedure
    tsvector_update_trigger(msg_content_tsy, 'public.chinese', msg_content);

-- 群通知
create table group_notify
(
    id bigserial primary key, -- 群通知 id
    relation_id bigint references relations (id) on delete cascade on update cascade, -- 关系 id（外键）
    msg_content text not null, -- 消息内容
    msg_expand json, -- 消息扩展信息
    account_id bigint references accounts (id) on delete cascade on update cascade, -- 发送账号 id（外键）
    create_at timestamptz not null default now(), -- 创建时间
    read_ids bigint[] not null default '{}'::bigint[], -- 已读用户 id 集合
    msg_content_tsv tsvector -- 消息分词
);

-- 分词索引
create index group_notify_msg_content_tsv on group_notify using gin (to_tsvector('chinese', msg_content));
-- 触发器更新 group_notify_msg_content_tsv
create trigger group_notify_msg_content_tsv
    before insert or update on group_notify
    for each row
execute procedure tsvector_update_trigger(msg_content_tsv, 'public.chinese', msg_content);

-- --
-- 创建更新 pin 时间戳的函数
create
    or replace function pin_timestamp() returns trigger as
$$
begin
    if
        new.is_pin then
        new.pin_time = now();
    end if;
    return new;
end;
$$
    language plpgsql;

-- 更新关系设置 pin 时间戳触发器
create trigger pin_timestamp_relations_settings_trigger
    before update of is_pin
    on settings
    for each row
execute procedure pin_timestamp();

-- 更新消息 pin 时间戳触发器
create trigger pin_timestamp_messages_trigger
    before update of is_pin
    on messages
    for each row
execute procedure pin_timestamp();

-- --
-- 创建更新时间戳的函数
create
    or replace function cs_timestamp() returns trigger as
$$
begin
    new.update_at = now();
    return new;
end;
$$
    language plpgsql;

-- 申请表更新时间戳触发器
create trigger application_update_at_trigger
    before update
    on applications
    for each row
execute procedure cs_timestamp();

-- --
-- 创建更新 show 时间戳的函数
create or replace function show_timestamp() returns trigger as
$$
begin
    if
        new.is_show then
        new.last_show = now();
    end if;
    return new;
end;
$$
    language plpgsql;

-- 更新关系设置 show 时间戳触发器
create trigger show_timestamp_trigger
    before update of is_show
    on settings
    for each row
execute procedure show_timestamp();