-- +goose Up

ALTER TABLE "public"."chat_ai_unknown_issue_stats"
    ADD COLUMN "sample_openid"      varchar(80) NOT NULL DEFAULT '',
    ADD COLUMN "sample_rel_user_id" int4        NOT NULL DEFAULT 0,
    ADD COLUMN "sample_dialogue_id" int4        NOT NULL DEFAULT 0,
    ADD COLUMN "sample_session_id"  int4        NOT NULL DEFAULT 0,
    ADD COLUMN "sample_message_id"  int4        NOT NULL DEFAULT 0,
    ADD COLUMN "last_dialogue_id"   int4        NOT NULL DEFAULT 0,
    ADD COLUMN "last_session_id"    int4        NOT NULL DEFAULT 0,
    ADD COLUMN "last_message_id"    int4        NOT NULL DEFAULT 0,
    ADD COLUMN "last_trigger_time"  int4        NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."sample_openid" IS '样本会话openid';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."sample_rel_user_id" IS '样本会话关联用户ID';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."sample_dialogue_id" IS '样本会话dialogue_id';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."sample_session_id" IS '样本会话session_id';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."sample_message_id" IS '样本会话消息ID';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."last_dialogue_id" IS '最近会话dialogue_id';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."last_session_id" IS '最近会话session_id';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."last_message_id" IS '最近会话消息ID';
COMMENT ON COLUMN "public"."chat_ai_unknown_issue_stats"."last_trigger_time" IS '最近触发时间';

