-- +goose Up

CREATE TABLE "chat_ai_library"
(
    "id"            serial        NOT NULL primary key,
    "admin_user_id" int4          NOT NULL DEFAULT 0,
    "library_name"  varchar(100)  NOT NULL DEFAULT '',
    "library_intro" varchar(1000) NOT NULL DEFAULT '',
    "model_config_id" int4         NOT NULL DEFAULT 0,
    "use_model"       varchar(100) NOT NULL DEFAULT '',
    "avatar"  varchar(512) NOT NULL DEFAULT '',
    "create_time"   int4          NOT NULL DEFAULT 0,
    "update_time"   int4          NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_library" ("admin_user_id");
CREATE INDEX ON "chat_ai_library" ("model_config_id");

COMMENT ON TABLE "chat_ai_library" IS '文档问答机器人-知识库';

COMMENT ON COLUMN "chat_ai_library"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_library"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_library"."library_name" IS '知识库名称';
COMMENT ON COLUMN "chat_ai_library"."library_intro" IS '知识库简介';
COMMENT ON COLUMN "chat_ai_library"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_library"."update_time" IS '更新时间';
COMMENT ON COLUMN "chat_ai_library"."model_config_id" IS '模型配置ID';
COMMENT ON COLUMN "chat_ai_library"."use_model" IS '使用的模型(枚举值)';
COMMENT ON COLUMN "chat_ai_library"."avatar" IS '知识库头像';
