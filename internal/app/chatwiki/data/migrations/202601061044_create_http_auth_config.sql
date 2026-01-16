-- +goose Up

CREATE TABLE "public"."http_node_auth_config" (
  "id"                bigserial         NOT NULL PRIMARY KEY,
  "staff_user_id"     int4              NOT NULL DEFAULT 0,
  "admin_user_id"     int4              NOT NULL DEFAULT 0,
  "create_time"       int4              NOT NULL DEFAULT 0,
  "update_time"       int4              NOT NULL DEFAULT 0,
  "auth_key"          varchar(30)       NOT NULL DEFAULT '',
  "auth_value"        varchar(255)      NOT NULL DEFAULT '',
  "auth_value_addto"  varchar(30)       NOT NULL DEFAULT '',
  "auth_remark"       varchar(255)      NOT NULL DEFAULT ''
); 

CREATE INDEX ON "public"."http_node_auth_config" ("staff_user_id");

COMMENT ON COLUMN "public"."http_node_auth_config"."staff_user_id" IS '客服user_id';

COMMENT ON COLUMN "public"."http_node_auth_config"."admin_user_id" IS '管理员user_id';

COMMENT ON COLUMN "public"."http_node_auth_config"."create_time" IS '创建时间';

COMMENT ON COLUMN "public"."http_node_auth_config"."update_time" IS '更新时间';

COMMENT ON COLUMN "public"."http_node_auth_config"."auth_key" IS '鉴权key';

COMMENT ON COLUMN "public"."http_node_auth_config"."auth_value" IS '鉴权value';

COMMENT ON COLUMN "public"."http_node_auth_config"."auth_value_addto" IS '加入到HEADERS、PARAMS、BODY';

COMMENT ON COLUMN "public"."http_node_auth_config"."auth_remark" IS '备注';

COMMENT ON TABLE "public"."http_node_auth_config" IS 'http节点鉴权配置表';