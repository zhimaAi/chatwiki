-- +goose Up

create table "public"."wechat_official_article_sync"
(
    "id"                        serial      not null primary key,
    "create_time"               int4        NOT NULL DEFAULT 0,
    "update_time"               int4        NOT NULL DEFAULT 0,
    "admin_user_id"             int4        not null default 0,
    "app_id"                    varchar(50) not null default '',
    "sync_type"                 int4        not null default 0,
    "sync_comment_switch"       int4        not null default 0,
    "ai_comment_switch"         int4        not null default 0,
    "replay_his_comment_switch" int4        not null default 1
);

CREATE INDEX on "wechat_official_article_sync"("admin_user_id");
CREATE INDEX on "wechat_official_article_sync"("app_id");

COMMENT ON TABLE "public"."wechat_official_article_sync" IS '公众号文章同步记录';
COMMENT ON COLUMN "public"."wechat_official_article_sync"."id" IS 'ID';
COMMENT ON COLUMN "public"."wechat_official_article_sync"."admin_user_id" IS '管理员ID';
COMMENT ON COLUMN "public"."wechat_official_article_sync"."app_id" IS '应用ID';
COMMENT ON COLUMN "public"."wechat_official_article_sync"."sync_type" IS '同步时间类型';
COMMENT ON COLUMN "public"."wechat_official_article_sync"."sync_comment_switch" IS '同步评论开关，0:关闭，1:开启';
COMMENT ON COLUMN "public"."wechat_official_article_sync"."ai_comment_switch" IS 'AI评论开关，0:关闭，1:开启';
COMMENT ON COLUMN "public"."wechat_official_article_sync"."replay_his_comment_switch" IS '历史评论是否处理，2:不处理，1:处理';
