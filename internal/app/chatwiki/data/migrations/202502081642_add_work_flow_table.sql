-- +goose Up

ALTER TABLE "chat_ai_robot"
    ADD COLUMN "start_node_key" varchar(32) NOT NULL DEFAULT '';
COMMENT ON COLUMN "chat_ai_robot"."start_node_key" IS '开始节点key';

CREATE TABLE "work_flow_node"
(
    "id"             bigserial   NOT NULL primary key,
    "admin_user_id"  int4        NOT NULL DEFAULT 0,
    "data_type"      int2        NOT NULL DEFAULT 0,
    "robot_id"       int4        NOT NULL DEFAULT 0,
    "node_type"      int4        NOT NULL DEFAULT 0,
    "node_name"      varchar(32) NOT NULL DEFAULT '',
    "node_key"       varchar(32) NOT NULL DEFAULT '',
    "node_params"    text        NOT NULL DEFAULT '{}',
    "node_info_json" text        NOT NULL DEFAULT '{}',
    "next_node_key"  varchar(32) NOT NULL DEFAULT '',
    "create_time"    int4        NOT NULL DEFAULT 0,
    "update_time"    int4        NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "work_flow_node" ("node_key", "robot_id", "data_type");
CREATE INDEX ON "work_flow_node" ("admin_user_id");
CREATE INDEX ON "work_flow_node" ("robot_id");

COMMENT ON COLUMN "work_flow_node"."id" IS 'ID';
COMMENT ON COLUMN "work_flow_node"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "work_flow_node"."data_type" IS '数据类型:1草稿,2发布';
COMMENT ON COLUMN "work_flow_node"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "work_flow_node"."node_type" IS '节点类型:0图的edges,其他的看代码';
COMMENT ON COLUMN "work_flow_node"."node_name" IS '节点名称';
COMMENT ON COLUMN "work_flow_node"."node_key" IS '节点key(hash)';
COMMENT ON COLUMN "work_flow_node"."node_params" IS '节点配置参数json';
COMMENT ON COLUMN "work_flow_node"."node_info_json" IS '图的单节点信息json,前端传啥存啥';
COMMENT ON COLUMN "work_flow_node"."next_node_key" IS '下一个节点key';
COMMENT ON COLUMN "work_flow_node"."create_time" IS '创建时间';
COMMENT ON COLUMN "work_flow_node"."update_time" IS '更新时间';

COMMENT ON TABLE "work_flow_node" IS '工作流-节点列表';


CREATE TABLE "work_flow_logs"
(
    "id"            bigserial    NOT NULL primary key,
    "admin_user_id" int4         NOT NULL DEFAULT 0,
    "robot_id"      int4         NOT NULL DEFAULT 0,
    "openid"        varchar(100) NOT NULL DEFAULT '',
    "run_node_keys" text         NOT NULL DEFAULT '',
    "run_logs"      text         NOT NULL DEFAULT '{}',
    "create_time"   int4         NOT NULL DEFAULT 0,
    "update_time"   int4         NOT NULL DEFAULT 0
);

CREATE INDEX ON "work_flow_logs" ("admin_user_id");
CREATE INDEX ON "work_flow_logs" ("robot_id");
CREATE INDEX ON "work_flow_logs" ("openid");

COMMENT ON COLUMN "work_flow_logs"."id" IS 'ID';
COMMENT ON COLUMN "work_flow_logs"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "work_flow_logs"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "work_flow_logs"."openid" IS '客户openid';
COMMENT ON COLUMN "work_flow_logs"."run_node_keys" IS '运行的节点key集';
COMMENT ON COLUMN "work_flow_logs"."run_logs" IS '运行日志json';
COMMENT ON COLUMN "work_flow_logs"."create_time" IS '创建时间';
COMMENT ON COLUMN "work_flow_logs"."update_time" IS '更新时间';

COMMENT ON TABLE "work_flow_logs" IS '工作流-运行日志';