-- +goose Up

CREATE TABLE IF NOT EXISTS "chat_ai_robot_multilingual_config"
(
    "id"                         serial8       NOT NULL,
    "admin_user_id"              int8          NOT NULL DEFAULT 0,
    "robot_id"                   int8          NOT NULL DEFAULT 0,
    "lang_key"                   varchar(20)   NOT NULL DEFAULT '',
    "welcomes"                   text          NOT NULL DEFAULT '',
    "unknown_question_prompt"    jsonb         NOT NULL DEFAULT '{"content":"","question":[]}',
    "tips_before_answer_switch"  bool          NOT NULL DEFAULT true,
    "tips_before_answer_content" text          NOT NULL DEFAULT '',
    "enable_common_question"     bool          NOT NULL DEFAULT false,
    "common_question_list"       jsonb         NOT NULL DEFAULT '[]',
    "create_time"                int8          NOT NULL DEFAULT 0,
    "update_time"                int8          NOT NULL DEFAULT 0,
    PRIMARY KEY ("id"),
    UNIQUE ("robot_id", "lang_key")
);

CREATE INDEX IF NOT EXISTS "idx_chat_ai_robot_multilingual_config_admin_user_id"
    ON "chat_ai_robot_multilingual_config" ("admin_user_id");
CREATE INDEX IF NOT EXISTS "idx_chat_ai_robot_multilingual_config_robot_id"
    ON "chat_ai_robot_multilingual_config" ("robot_id");

COMMENT ON TABLE "chat_ai_robot_multilingual_config" IS '机器人多语言基础配置';
COMMENT ON COLUMN "chat_ai_robot_multilingual_config"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_robot_multilingual_config"."robot_id" IS '机器人ID';
COMMENT ON COLUMN "chat_ai_robot_multilingual_config"."lang_key" IS '语言key，如ch/en/ja';
COMMENT ON COLUMN "chat_ai_robot_multilingual_config"."welcomes" IS '欢迎语JSON';
COMMENT ON COLUMN "chat_ai_robot_multilingual_config"."unknown_question_prompt" IS '未知问题提示语JSON';
COMMENT ON COLUMN "chat_ai_robot_multilingual_config"."tips_before_answer_switch" IS '答案生成中提示语开关';
COMMENT ON COLUMN "chat_ai_robot_multilingual_config"."tips_before_answer_content" IS '答案生成中提示语';
COMMENT ON COLUMN "chat_ai_robot_multilingual_config"."enable_common_question" IS '常见问题开关';
COMMENT ON COLUMN "chat_ai_robot_multilingual_config"."common_question_list" IS '常见问题列表';
