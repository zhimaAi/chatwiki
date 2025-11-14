-- +goose Up

CREATE TABLE "public"."mcp_server" (
    "id" serial NOT NULL primary key,
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "name" varchar(100) NOT NULL DEFAULT '',
    "description" varchar(500) NOT NULL DEFAULT '',
    "avatar" varchar(500) NOT NULL DEFAULT '',
    "auth_type" int4 NOT NULL DEFAULT 1,
    "api_key" varchar(500) NOT NULL DEFAULT '',
    "publish_status" int2 NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."mcp_server" ("admin_user_id");
CREATE INDEX ON "public"."mcp_server" ("api_key");

COMMENT ON TABLE "public"."mcp_server" IS 'MCP服务器';
COMMENT ON COLUMN "public"."mcp_server"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."mcp_server"."name" IS 'MCP服务器名称';
COMMENT ON COLUMN "public"."mcp_server"."description" IS 'MCP服务器描述';
COMMENT ON COLUMN "public"."mcp_server"."avatar" IS '机器人头像';
COMMENT ON COLUMN "public"."mcp_server"."auth_type" IS '授权方式 1 API KEY';
COMMENT ON COLUMN "public"."mcp_server"."api_key" IS 'API KEY';
COMMENT ON COLUMN "public"."mcp_server"."publish_status" IS '发布状态 1 未发布 2 已发布';