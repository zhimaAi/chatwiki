-- +goose Up
CREATE TABLE "public"."wechat_official_account_batch_send_task"
(
    "id"                     serial                                      NOT NULL primary key,
    "admin_user_id"          INT4                                        NOT NULL DEFAULT 0,
    "create_time"            INT4                                        NOT NULL DEFAULT 0,
    "update_time"            INT4                                        NOT NULL DEFAULT 0,
    "app_id"                 VARCHAR(30) COLLATE "pg_catalog"."default"  NOT NULL DEFAULT '' :: CHARACTER VARYING,
    "draft_id"               INT4                                        NOT NULL DEFAULT 0,
    "comment_status"         INT2                                        NOT NULL DEFAULT 0,
    "open_status"            INT2                                        NOT NULL DEFAULT 0,
    "send_time"              INT4                                        NOT NULL DEFAULT 0,
    "is_top"                 INT2                                        NOT NULL DEFAULT 0,
    "to_user_type"           INT2                                        NOT NULL DEFAULT 0,
    "to_user"                TEXT COLLATE "pg_catalog"."default"         NOT NULL DEFAULT '' :: TEXT,
    "res_total"              INT2                                        NOT NULL DEFAULT 0,
    "res_succ"               INT2                                        NOT NULL DEFAULT 0,
    "res_err"                INT2                                        NOT NULL DEFAULT 0,
    "msg_data_id"            VARCHAR(50) COLLATE "pg_catalog"."default"  NOT NULL DEFAULT '' :: CHARACTER VARYING,
    "send_status"            INT2                                        NOT NULL DEFAULT 0,
    "common_total"           INT2                                        NOT NULL DEFAULT 0,
    "access_key"             VARCHAR(50) COLLATE "pg_catalog"."default"  NOT NULL DEFAULT '' :: CHARACTER VARYING,
    "comment_rule_id"        INT2                                        NOT NULL DEFAULT 0,
    "task_name"              VARCHAR(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT '' :: CHARACTER VARYING,
    "send_res"               TEXT COLLATE "pg_catalog"."default"         NOT NULL DEFAULT '' :: TEXT,
    "last_comment_sync_time" INT4                                        NOT NULL DEFAULT 0,
    "send_msg_id"            VARCHAR(50) COLLATE "pg_catalog"."default"  NOT NULL DEFAULT '' :: CHARACTER VARYING,
    "max_comment_id"         INT2                                        NOT NULL DEFAULT 0,
    "ai_comment_status"      INT2                                        NOT NULL DEFAULT 0
);
ALTER TABLE "public"."wechat_official_account_batch_send_task" OWNER TO "chatwiki";
CREATE INDEX "wechat_official_send_task_" ON "public"."wechat_official_account_batch_send_task" USING btree ( "admin_user_id" "pg_catalog"."int4_ops" ASC NULLS LAST, "app_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST, "draft_id" "pg_catalog"."int4_ops" ASC NULLS LAST );
CREATE INDEX "wechat_official_send_task_access_key" ON "public"."wechat_official_account_batch_send_task" USING btree ( "admin_user_id" "pg_catalog"."int4_ops" ASC NULLS LAST, "access_key" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST );
CREATE INDEX "wechat_official_send_task_send_status_time" ON "public"."wechat_official_account_batch_send_task" USING btree ( "send_status" "pg_catalog"."int2_ops" ASC NULLS LAST, "send_time" "pg_catalog"."int4_ops" ASC NULLS LAST );
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."id" IS 'ID';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."update_time" IS '更新时间';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."app_id" IS '公众号ID';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."draft_id" IS '草稿ID';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."comment_status" IS '评论状态,0:关闭，1:开启';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."open_status" IS '任务状态 0:关闭，1:开启';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."send_time" IS '发送时间，0:立即群发，非0 发送时间';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."is_top" IS '是否置顶，0:不是，1:是';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."to_user_type" IS '发送对象类型，0:全部粉丝，1:标签粉丝，2:openid';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."to_user" IS 'to_user_type 非0时存在';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."res_total" IS '发送数';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."res_succ" IS '成功数';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."res_err" IS '失败数';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."msg_data_id" IS '群发消息ID';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."send_status" IS '群发状态，0:未开始，1:执行中，2:已完成，3:执行失败-1:已删除';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."common_total" IS '评论数';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."access_key" IS '应用关联机器人key';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."comment_rule_id" IS '评论规则ID，0为默认评论';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."task_name" IS '群发规则名';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."send_res" IS '微信群发执行结果';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."last_comment_sync_time" IS '上次评论同步时间';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."send_msg_id" IS '群发消息ID';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."max_comment_id" IS '最大评论ID';
COMMENT ON COLUMN "public"."wechat_official_account_batch_send_task"."ai_comment_status" IS '是否启用AI评论';
COMMENT ON TABLE "public"."wechat_official_account_batch_send_task" IS '公众号文章群发任务';