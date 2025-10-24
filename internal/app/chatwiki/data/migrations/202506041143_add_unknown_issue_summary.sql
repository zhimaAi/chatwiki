-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "unknown_summary_status"          int2         NOT NULL DEFAULT 0,
    ADD COLUMN "unknown_summary_model_config_id" int4         NOT NULL DEFAULT 0,
    ADD COLUMN "unknown_summary_use_model"       varchar(100) NOT NULL DEFAULT '',
    ADD COLUMN "unknown_summary_similarity"      float4       NOT NULL DEFAULT 0.8;

COMMENT ON COLUMN "public"."chat_ai_robot"."unknown_summary_status" IS '未知问题总结开关状态';
COMMENT ON COLUMN "public"."chat_ai_robot"."unknown_summary_model_config_id" IS '未知问题总结模型配置ID';
COMMENT ON COLUMN "public"."chat_ai_robot"."unknown_summary_use_model" IS '未知问题总结使用的模型';
COMMENT ON COLUMN "public"."chat_ai_robot"."unknown_summary_similarity" IS '未知问题总结判定相似度';


CREATE TABLE "public"."chat_ai_unknown_issue_summary"
(
    "id"              serial         NOT NULL primary key,
    "admin_user_id"   int4           NOT NULL DEFAULT 0,
    "robot_id"        int4           NOT NULL DEFAULT 0,
    "trigger_day"     int4           NOT NULL DEFAULT 0,
    "question"        varchar(10000) NOT NULL DEFAULT '',
    "embedding"       vector,
    "unknown_list"    jsonb          NOT NULL DEFAULT '[]',
    "unknown_total"   int4           NOT NULL DEFAULT 0,
    "answer"          varchar(10000) NOT NULL DEFAULT '',
    "images"          json                    DEFAULT '[]',
    "to_library_id"   int4           NOT NULL DEFAULT 0,
    "to_library_name" varchar(100)   NOT NULL DEFAULT '',
    "create_time"     int4           NOT NULL DEFAULT 0,
    "update_time"     int4           NOT NULL DEFAULT 0
);

ALTER TABLE "public"."chat_ai_unknown_issue_summary"
    ADD CONSTRAINT check_images_array CHECK (json_typeof(images) = 'array');

CREATE INDEX ON "public"."chat_ai_unknown_issue_summary"
    ("admin_user_id", "robot_id", "trigger_day");
CREATE INDEX ON "public"."chat_ai_unknown_issue_summary" ("robot_id", "trigger_day");

COMMENT ON TABLE "public"."chat_ai_unknown_issue_summary" IS '未知问题总结';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."trigger_day" IS '日期(yyyymmdd)';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."question" IS '聚类问题';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."embedding" IS '可变维度向量';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."unknown_list" IS '未知问题json';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."unknown_total" IS '未知问题数量';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."answer" IS '设置的答案';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."images" IS '图片列表';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."to_library_id" IS '导入的知识库ID';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."to_library_name" IS '导入的知识库名称';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_summary"."update_time" IS '更新时间';