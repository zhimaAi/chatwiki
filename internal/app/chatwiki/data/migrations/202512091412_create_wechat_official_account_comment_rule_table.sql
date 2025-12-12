-- +goose Up
CREATE TABLE "public"."wechat_official_account_comment_rule"
(
    "id"                    serial                                      NOT NULL primary key,
    "admin_user_id"         int4                                        NOT NULL DEFAULT 0,
    "create_time"           int4                                        NOT NULL DEFAULT 0,
    "update_time"           int4                                        NOT NULL DEFAULT 0,
    "rule_name"             varchar(50) COLLATE "pg_catalog"."default"  NOT NULL DEFAULT ''::character varying,
    "use_model"             varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
    "delete_comment_switch" int2                                        NOT NULL DEFAULT 0,
    "delete_comment_rule"   jsonb,
    "reply_comment_switch"  int2                                        NOT NULL DEFAULT 0,
    "reply_comment_rule"    jsonb,
    "elect_comment_switch"  int2                                        NOT NULL DEFAULT 0,
    "elect_comment_rule"    jsonb,
    "is_default"            int2                                        NOT NULL DEFAULT 0,
    "switch"                int2                                        NOT NULL DEFAULT 0,
    "model_config_id"       int2                                        NOT NULL DEFAULT 0
);

ALTER TABLE "public"."wechat_official_account_comment_rule"
    OWNER TO "chatwiki";

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."id" IS 'ID';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."admin_user_id" IS '管理员用户ID';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."create_time" IS '创建时间';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."update_time" IS '更新时间';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."rule_name" IS '规则名称';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."use_model" IS '使用的模型(枚举值)';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."delete_comment_switch" IS '删除评论开关，0:关闭，1:开启';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."delete_comment_rule" IS '删除规则';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."reply_comment_switch" IS '自动回复开关。0:关闭。1:开启';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."reply_comment_rule" IS '自动回复规则';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."elect_comment_switch" IS '自动精选开关，0:关闭，1:开启';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."elect_comment_rule" IS '自动精选规则';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."is_default" IS '是否为默认规则，0:否，1:是';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."switch" IS '规则是否开启';

COMMENT ON COLUMN "public"."wechat_official_account_comment_rule"."model_config_id" IS '公众号文章AI评论规则';