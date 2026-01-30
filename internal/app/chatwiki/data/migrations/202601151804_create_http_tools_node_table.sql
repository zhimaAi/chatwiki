-- +goose Up

create table "public"."http_tools_node" (
    "id" serial not null primary key,
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0,
    "admin_user_id" int4 not null default 0,
    "http_tool_id" int4 not null default 0,
    "node_key" varchar(100) not null default '',
    "node_name" varchar(100) not null default '',
    "node_name_en" varchar(100) not null default '',
    "node_description" varchar(1000) not null default '',
    "node_remark" varchar(500) not null default '',
    "data_raw" TEXT COLLATE "pg_catalog"."default"
);

-- 添加索引 node_key
CREATE INDEX on "http_tools_node"("admin_user_id");
CREATE INDEX on "http_tools_node"("http_tool_id");
CREATE INDEX on "http_tools_node"("node_key");

COMMENT ON TABLE "public"."http_tools_node" IS 'http工具节点';
COMMENT ON COLUMN "public"."http_tools_node"."id" IS 'ID';
COMMENT ON COLUMN "public"."http_tools_node"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."http_tools_node"."update_time" IS '更新时间';
COMMENT ON COLUMN "public"."http_tools_node"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."http_tools_node"."http_tool_id" IS 'http工具ID';
COMMENT ON COLUMN "public"."http_tools_node"."node_key" IS 'http工具节点key（唯一）删除检测';
COMMENT ON COLUMN "public"."http_tools_node"."node_name" IS '名称';
COMMENT ON COLUMN "public"."http_tools_node"."node_name_en" IS '工具英文名称';
COMMENT ON COLUMN "public"."http_tools_node"."node_description" IS '描述';
COMMENT ON COLUMN "public"."http_tools_node"."node_remark" IS '备注';
COMMENT ON COLUMN "public"."http_tools_node"."data_raw" IS '原始数据';
