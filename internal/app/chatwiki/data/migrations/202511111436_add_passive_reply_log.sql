-- +goose Up

CREATE TABLE "public"."chat_ai_passive_reply_log"
(
    "id"            bigserial    NOT NULL primary key,
    "admin_user_id" int4         NOT NULL DEFAULT 0,
    "app_id"        varchar(100) NOT NULL DEFAULT '',
    "msgid"         varchar(100) NOT NULL DEFAULT '',
    "openid"        varchar(100) NOT NULL DEFAULT '',
    "content"       text         NOT NULL DEFAULT '',
    "status"        int2         NOT NULL DEFAULT 0,
    "replys"        text         NOT NULL DEFAULT '[]',
    "images"        text         NOT NULL DEFAULT '[]',
    "create_time"   int4         NOT NULL DEFAULT 0,
    "update_time"   int4         NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_passive_reply_log" ("admin_user_id");
CREATE INDEX ON "public"."chat_ai_passive_reply_log" ("app_id");
CREATE INDEX ON "public"."chat_ai_passive_reply_log" ("msgid", "app_id");
CREATE INDEX ON "public"."chat_ai_passive_reply_log" ("openid");
CREATE INDEX ON "public"."chat_ai_passive_reply_log" ("status");

COMMENT ON TABLE "public"."chat_ai_passive_reply_log" IS '公众号被动回复日志';
COMMENT ON COLUMN "public"."chat_ai_passive_reply_log"."id" IS 'ID';
COMMENT ON COLUMN "public"."chat_ai_passive_reply_log"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."chat_ai_passive_reply_log"."app_id" IS '应用ID';
COMMENT ON COLUMN "public"."chat_ai_passive_reply_log"."msgid" IS '消息ID';
COMMENT ON COLUMN "public"."chat_ai_passive_reply_log"."openid" IS '客户ID';
COMMENT ON COLUMN "public"."chat_ai_passive_reply_log"."content" IS '提问内容';
COMMENT ON COLUMN "public"."chat_ai_passive_reply_log"."status" IS '状态:0生成中,1生成完毕';
COMMENT ON COLUMN "public"."chat_ai_passive_reply_log"."replys" IS '文本回复列表';
COMMENT ON COLUMN "public"."chat_ai_passive_reply_log"."images" IS '图片回复列表';
COMMENT ON COLUMN "public"."chat_ai_passive_reply_log"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."chat_ai_passive_reply_log"."update_time" IS '更新时间';
