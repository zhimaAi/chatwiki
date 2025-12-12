-- +goose Up
CREATE TABLE "public"."wechat_official_comment_list"
(
    "id"                     serial                                      NOT NULL PRIMARY KEY,
    "admin_user_id"          INT4                                        NOT NULL DEFAULT 0,
    "create_time"            INT4                                        NOT NULL DEFAULT 0,
    "update_time"            INT4                                        NOT NULL DEFAULT 0,
    "msg_data_id"            VARCHAR(50) COLLATE "pg_catalog"."default"  NOT NULL DEFAULT '' :: CHARACTER VARYING,
    "access_key"             VARCHAR(50) COLLATE "pg_catalog"."default"  NOT NULL DEFAULT '' :: CHARACTER VARYING,
    "comment_rule_id"        INT2                                        NOT NULL DEFAULT 0,
    "user_comment_id"        INT2                                        NOT NULL DEFAULT 0,
    "comment_create_time"    INT4                                        NOT NULL DEFAULT 0,
    "content_text"           TEXT COLLATE "pg_catalog"."default",
    "open_id"                VARCHAR(50) COLLATE "pg_catalog"."default"  NOT NULL DEFAULT '' :: CHARACTER VARYING,
    "reply_comment_text"     TEXT COLLATE "pg_catalog"."default"         NOT NULL DEFAULT '' :: TEXT,
    "reply_create_time"      INT4                                        NOT NULL DEFAULT 0,
    "comment_type"           INT2                                        NOT NULL DEFAULT 0,
    "ai_comment_rule_status" INT2                                        NOT NULL DEFAULT 0,
    "task_id"                INT4                                        NOT NULL DEFAULT 0,
    "draft_id"               INT2                                        NOT NULL DEFAULT 0,
    "ai_comment_rule_text"   VARCHAR(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '' :: CHARACTER VARYING,
    "ai_comment_result"      JSONB,
    "delete_status"          INT2                                        NOT NULL DEFAULT 0,
    "ai_exec_time"           INT4                                        NOT NULL DEFAULT 0
);
ALTER TABLE "public"."wechat_official_comment_list" OWNER TO "chatwiki";
CREATE INDEX "wechat_official_comment_list_access_msg_comment_index" ON "public"."wechat_official_comment_list" USING btree (
    "msg_data_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
    "access_key" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
    "user_comment_id" "pg_catalog"."int2_ops" ASC NULLS LAST
);
CREATE UNIQUE INDEX "wechat_official_comment_list_admin_msg_comment_id" ON "public"."wechat_official_comment_list" USING btree ( "admin_user_id" "pg_catalog"."int4_ops" ASC NULLS LAST, "msg_data_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST, "user_comment_id" "pg_catalog"."int2_ops" ASC NULLS LAST );
CREATE INDEX "wechat_official_comment_list_admin_task_id" ON "public"."wechat_official_comment_list" USING btree ( "admin_user_id" "pg_catalog"."int4_ops" ASC NULLS LAST, "task_id" "pg_catalog"."int4_ops" ASC NULLS LAST );
COMMENT ON COLUMN "public"."wechat_official_comment_list"."id" IS 'ID';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."update_time" IS '更新时间';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."msg_data_id" IS '群发消息ID';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."access_key" IS '应用关联机器人key';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."comment_rule_id" IS '评论规则ID，0为默认评论';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."user_comment_id" IS '评论ID';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."comment_create_time" IS '评论时间';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."content_text" IS '评论内容';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."open_id" IS '评论人openid';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."reply_comment_text" IS '回复评论内容';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."reply_create_time" IS '回复时间';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."comment_type" IS '评论类型';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."ai_comment_rule_status" IS '是否经过AI评论规则执行，0:未进行，1:执行过';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."task_id" IS '群发任务ID';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."draft_id" IS '草稿ID';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."ai_comment_rule_text" IS '触发条件';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."ai_comment_result" IS '处理结果';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."delete_status" IS '删除状态，0:未删除，1:已删除';
COMMENT ON COLUMN "public"."wechat_official_comment_list"."ai_exec_time" IS 'AI精选执行时间';
COMMENT ON TABLE "public"."wechat_official_comment_list" IS '公众号文章评论';