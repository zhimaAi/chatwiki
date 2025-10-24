-- +goose Up

CREATE TABLE "work_flow_node_version"
(
    "id"             bigserial   NOT NULL primary key,
    "admin_user_id"  int4        NOT NULL DEFAULT 0,
    "robot_id"       int4        NOT NULL DEFAULT 0,
    "work_flow_version_id"       int4        NOT NULL DEFAULT 0,
    "node_type"      int4        NOT NULL DEFAULT 0,
    "node_name"      varchar(32) NOT NULL DEFAULT '',
    "node_key"       varchar(32) NOT NULL DEFAULT '',
    "node_params"    text        NOT NULL DEFAULT '{}',
    "node_info_json" text        NOT NULL DEFAULT '{}',
    "next_node_key"  varchar(32) NOT NULL DEFAULT '',
    "create_time"    int4        NOT NULL DEFAULT 0,
    "update_time"    int4        NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "work_flow_node_version" ("node_key", "robot_id", "work_flow_version_id");
CREATE INDEX ON "work_flow_node_version" ("admin_user_id");
CREATE INDEX ON "work_flow_node_version" ("robot_id");
CREATE INDEX ON "work_flow_node_version" ("work_flow_version_id");

COMMENT ON COLUMN "work_flow_node_version"."id" IS 'ID';
COMMENT ON COLUMN "work_flow_node_version"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "work_flow_node_version"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "work_flow_node_version"."work_flow_version_id" IS '机器人版本ID';
COMMENT ON COLUMN "work_flow_node_version"."node_type" IS '节点类型:0图的edges,其他的看代码';
COMMENT ON COLUMN "work_flow_node_version"."node_name" IS '节点名称';
COMMENT ON COLUMN "work_flow_node_version"."node_key" IS '节点key(hash)';
COMMENT ON COLUMN "work_flow_node_version"."node_params" IS '节点配置参数json';
COMMENT ON COLUMN "work_flow_node_version"."node_info_json" IS '图的单节点信息json,前端传啥存啥';
COMMENT ON COLUMN "work_flow_node_version"."next_node_key" IS '下一个节点key';
COMMENT ON COLUMN "work_flow_node_version"."create_time" IS '创建时间';
COMMENT ON COLUMN "work_flow_node_version"."update_time" IS '更新时间';

COMMENT ON TABLE "work_flow_node_version" IS '工作流-历史节点列表';



CREATE TABLE "work_flow_version"
(
    "id"             bigserial   NOT NULL primary key,
    "admin_user_id"  int4        NOT NULL DEFAULT 0,
    "robot_id"       int4        NOT NULL DEFAULT 0,
    "version"        varchar(32) NOT NULL DEFAULT '',
    "version_desc"           varchar(500) NOT NULL DEFAULT '',
    "create_time"    int4        NOT NULL DEFAULT 0,
    "update_time"    int4        NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "work_flow_version" ("robot_id", "version");
CREATE INDEX ON "work_flow_version" ("admin_user_id");


COMMENT ON COLUMN "work_flow_version"."id" IS 'ID';
COMMENT ON COLUMN "work_flow_version"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "work_flow_version"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "work_flow_version"."version" IS '版本号';
COMMENT ON COLUMN "work_flow_version"."version_desc" IS '版本描述';
COMMENT ON COLUMN "work_flow_version"."create_time" IS '创建时间';
COMMENT ON COLUMN "work_flow_version"."update_time" IS '更新时间';

COMMENT ON TABLE "work_flow_version" IS '工作流-历史版本';

