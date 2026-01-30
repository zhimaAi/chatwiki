-- +goose Up

create table "public"."http_tools" (
   "id" serial not null primary key,
   "create_time" int4 NOT NULL DEFAULT 0,
   "update_time" int4 NOT NULL DEFAULT 0,
   "admin_user_id" int4 not null default 0,
   "name" varchar(100) not null default '',
   "name_en" varchar(100) not null default '',
   "tool_key" varchar(100) not null default '',
   "avatar" varchar(500) not null default '',
   "description" varchar(1000) not null default ''
);

-- 添加索引 tool_key
CREATE INDEX on "public"."http_tools"("admin_user_id");
CREATE INDEX on "public"."http_tools"("tool_key");

COMMENT ON TABLE "public"."http_tools" IS 'http工具集';
COMMENT ON COLUMN "public"."http_tools"."id" IS 'ID';
COMMENT ON COLUMN "public"."http_tools"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."http_tools"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."http_tools"."update_time" IS '更新时间';
COMMENT ON COLUMN "public"."http_tools"."name" IS '名称';
COMMENT ON COLUMN "public"."http_tools"."name_en" IS '工具集英文名称';
COMMENT ON COLUMN "public"."http_tools"."tool_key" IS '工具key（唯一）添加检测';
COMMENT ON COLUMN "public"."http_tools"."avatar" IS '图标';
COMMENT ON COLUMN "public"."http_tools"."description" IS '描述';