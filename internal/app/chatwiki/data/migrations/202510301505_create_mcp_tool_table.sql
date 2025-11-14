-- +goose Up

CREATE TABLE "public"."mcp_tool" (
    "id" serial NOT NULL primary key,
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "server_id" int4 NOT NULL DEFAULT 0,
    "robot_id" int4 NOT NULL DEFAULT 0,
    "name" varchar(255) NOT NULL DEFAULT ''
);

CREATE INDEX ON "public"."mcp_tool" ("admin_user_id", "server_id");

COMMENT ON TABLE "public"."mcp_tool" IS 'MCP工具';
COMMENT ON COLUMN "public"."mcp_tool"."admin_user_id" IS '管理员ID';
COMMENT ON COLUMN "public"."mcp_tool"."server_id" IS 'MCP服务器ID';
COMMENT ON COLUMN "public"."mcp_tool"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."mcp_tool"."name" IS '工具名称';
