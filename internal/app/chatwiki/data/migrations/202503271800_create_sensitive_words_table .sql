-- +goose Up

CREATE TABLE "public"."sensitive_words"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "words" text          NOT NULL DEFAULT '',
    "trigger_type"      int2          NOT NULL DEFAULT 0,
    "creator" int4          NOT NULL DEFAULT 0,
    "status" int2          NOT NULL DEFAULT 0,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."sensitive_words" ("admin_user_id", "trigger_type");

COMMENT ON TABLE "public"."sensitive_words" IS '敏感词表';

COMMENT ON COLUMN "public"."sensitive_words"."id" IS 'ID';
COMMENT ON COLUMN "public"."sensitive_words"."admin_user_id" IS '用户ID';
COMMENT ON COLUMN "public"."sensitive_words"."words" IS '敏感词';
COMMENT ON COLUMN "public"."sensitive_words"."trigger_type" IS '触发类型 0:所有 1:指定';
COMMENT ON COLUMN "public"."sensitive_words"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."sensitive_words"."update_time" IS '更新时间';

CREATE TABLE "public"."sensitive_words_relation"
(
    "id"            serial        NOT NULL primary key,
    "words_id" int4          NOT NULL DEFAULT 0,
    "robot_id"      int2          NOT NULL DEFAULT 0,
    "status" int2          NOT NULL DEFAULT 0,
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."sensitive_words_relation" ("words_id");
CREATE INDEX ON "public"."sensitive_words_relation" ("robot_id");

COMMENT ON TABLE "public"."sensitive_words_relation" IS '敏感词表';

COMMENT ON COLUMN "public"."sensitive_words_relation"."id" IS 'ID';
COMMENT ON COLUMN "public"."sensitive_words_relation"."words_id" IS '敏感词id';
COMMENT ON COLUMN "public"."sensitive_words_relation"."robot_id" IS '机器人id';
COMMENT ON COLUMN "public"."sensitive_words_relation"."status" IS '1:启用 2:禁用';
COMMENT ON COLUMN "public"."sensitive_words_relation"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."sensitive_words_relation"."update_time" IS '更新时间';

