-- +goose Up

CREATE TABLE "public"."chat_ai_stat_library_robot_tip" (
    "id" serial NOT NULL primary key,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "library_id" int4 NOT NULL DEFAULT 0,
    "robot_id" int4 NOT NULL DEFAULT 0,
    "date_ymd" varchar(10) NOT NULL DEFAULT '',
    "tip" int4 NOT NULL DEFAULT 0,
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_stat_library_robot_tip" ("admin_user_id" , "date_ymd");
CREATE UNIQUE INDEX ON "public"."chat_ai_stat_library_robot_tip" ("admin_user_id" , "library_id" , "date_ymd" , "robot_id");

COMMENT ON TABLE "public"."chat_ai_stat_library_robot_tip" IS '知识库AI触发次数统计表';

COMMENT ON COLUMN "public"."chat_ai_stat_library_robot_tip"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_stat_library_robot_tip"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_stat_library_robot_tip"."library_id" IS '知识库ID';
COMMENT ON COLUMN "public"."chat_ai_stat_library_robot_tip"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."chat_ai_stat_library_robot_tip"."date_ymd" IS '日期：20251010';
COMMENT ON COLUMN "public"."chat_ai_stat_library_robot_tip"."tip" IS '触发次数';
COMMENT ON COLUMN "public"."chat_ai_stat_library_robot_tip"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_stat_library_robot_tip"."update_time" IS '更新时间';



CREATE TABLE "public"."chat_ai_stat_library_data_robot_tip" (
     "id" serial NOT NULL primary key,
     "admin_user_id" int4 NOT NULL DEFAULT 0,
     "robot_id" int4 NOT NULL DEFAULT 0,
     "library_id" int4 NOT NULL DEFAULT 0,
     "file_id" varchar(10) NOT NULL DEFAULT '',
     "data_id" int4 NOT NULL DEFAULT 0,
     "date_ymd" varchar(10) NOT NULL DEFAULT '',
     "tip" int4 NOT NULL DEFAULT 0,
     "create_time" int4 NOT NULL DEFAULT 0,
     "update_time" int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_stat_library_data_robot_tip" ("admin_user_id" , "date_ymd" , "library_id");
CREATE UNIQUE INDEX ON "public"."chat_ai_stat_library_data_robot_tip" ("admin_user_id" , "data_id" , "date_ymd" , "robot_id");

COMMENT ON TABLE "public"."chat_ai_stat_library_data_robot_tip" IS '知识库内容AI触发次数统计表';

COMMENT ON COLUMN "public"."chat_ai_stat_library_data_robot_tip"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_stat_library_data_robot_tip"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_stat_library_data_robot_tip"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."chat_ai_stat_library_data_robot_tip"."library_id" IS '知识库ID';
COMMENT ON COLUMN "public"."chat_ai_stat_library_data_robot_tip"."data_id" IS '知识库中文档ID';
COMMENT ON COLUMN "public"."chat_ai_stat_library_data_robot_tip"."data_id" IS '知识库数据ID';
COMMENT ON COLUMN "public"."chat_ai_stat_library_data_robot_tip"."date_ymd" IS '日期：20251010';
COMMENT ON COLUMN "public"."chat_ai_stat_library_data_robot_tip"."tip" IS '触发次数';
COMMENT ON COLUMN "public"."chat_ai_stat_library_data_robot_tip"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_stat_library_data_robot_tip"."update_time" IS '更新时间';