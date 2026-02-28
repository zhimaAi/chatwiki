-- +goose Up

CREATE TABLE "workbench_config"
(
    "id"                           bigserial       NOT NULL primary key,
    "admin_user_id"                int4            NOT NULL DEFAULT 0,
    "default_robot_id"             int4            NOT NULL DEFAULT 0,
    "enable_last_app_entry"        int2            NOT NULL DEFAULT 1,
    "create_time"                  int4            NOT NULL DEFAULT 0,
    "update_time"                  int4            NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "workbench_config" ("admin_user_id");

COMMENT ON TABLE "workbench_config" IS '工作台配置表';

COMMENT ON COLUMN "workbench_config"."id" IS 'ID';
COMMENT ON COLUMN "workbench_config"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "workbench_config"."default_robot_id" IS '首页默认显示的机器人ID';
COMMENT ON COLUMN "workbench_config"."enable_last_app_entry" IS '是否开启默认进入上一次访问的应用开关,1开启,0关闭';
COMMENT ON COLUMN "workbench_config"."create_time" IS '创建时间';
COMMENT ON COLUMN "workbench_config"."update_time" IS '更新时间';

CREATE TABLE "robot_history_visit"
(
    "id"                           bigserial       NOT NULL primary key,
    "admin_user_id"                int4            NOT NULL DEFAULT 0,
    "user_id"                      int4            NOT NULL DEFAULT 0,
    "robot_id"                     int4            NOT NULL DEFAULT 0,
    "create_time"                  int4            NOT NULL DEFAULT 0,
    "update_time"                  int4            NOT NULL DEFAULT 0
);

CREATE INDEX ON "robot_history_visit" ("admin_user_id","user_id");

COMMENT ON TABLE "robot_history_visit" IS '历史访问表';

COMMENT ON COLUMN "robot_history_visit"."id" IS 'ID';
COMMENT ON COLUMN "robot_history_visit"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "robot_history_visit"."user_id" IS '用户ID';
COMMENT ON COLUMN "robot_history_visit"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "robot_history_visit"."create_time" IS '创建时间';
COMMENT ON COLUMN "robot_history_visit"."update_time" IS '更新时间';


CREATE TABLE "workbench_top_robot"
(
    "id"                           bigserial       NOT NULL primary key,
    "admin_user_id"                int4            NOT NULL DEFAULT 0,
    "user_id"                      int4            NOT NULL DEFAULT 0,
    "robot_id"                     int4            NOT NULL DEFAULT 0,
    "create_time"                  int4            NOT NULL DEFAULT 0,
    "update_time"                  int4            NOT NULL DEFAULT 0
);

CREATE INDEX ON "workbench_top_robot" ("admin_user_id","user_id");

COMMENT ON TABLE "workbench_top_robot" IS '机器人置顶表';

COMMENT ON COLUMN "workbench_top_robot"."id" IS 'ID';
COMMENT ON COLUMN "workbench_top_robot"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "workbench_top_robot"."user_id" IS '用户ID';
COMMENT ON COLUMN "workbench_top_robot"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "workbench_top_robot"."create_time" IS '创建时间';
COMMENT ON COLUMN "workbench_top_robot"."update_time" IS '更新时间';

