-- +goose Up

create table "public"."mcp_provider" (
    "id" serial not null primary key,
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0,
    "admin_user_id" int4 not null default 0,
    "avatar" varchar(500) not null default '',
    "name" varchar(100) not null default '',
    "description" varchar(1000) not null default '',
    "url" varchar(500) not null default '',
    "request_timeout" int4 not null default 300,
    "client_type" int2 not null default 1,
    "has_auth" int2 not null default 0,
    "headers" jsonb not null default '[]',
    "tools" jsonb not null default '[]'
);

comment on table "public"."mcp_provider" is 'mcp提供商';
comment on column "public"."mcp_provider"."avatar" is '头像';
comment on column "public"."mcp_provider"."name" is '名称';
comment on column "public"."mcp_provider"."description" is '描述';
comment on column "public"."mcp_provider"."url" is '地址';
comment on column "public"."mcp_provider"."request_timeout" is '请求超时时间';
comment on column "public"."mcp_provider"."client_type" is '客户端类型 1,sse 2,sse';
comment on column "public"."mcp_provider"."has_auth" is '是否授权过 0未授权 1已授权';
comment on column "public"."mcp_provider"."headers" is '请求头';
comment on column "public"."mcp_provider"."tools" is '工具集';