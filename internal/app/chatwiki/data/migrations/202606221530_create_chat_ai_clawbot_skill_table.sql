-- +goose Up
CREATE TABLE "chat_ai_clawbot_user_skill"
(
    "id"               serial       NOT NULL primary key,
    "admin_user_id"    int4         NOT NULL DEFAULT 0,
    "skill_name"       varchar(50)  NOT NULL DEFAULT '',
    "remark_name"      varchar(50)  NOT NULL DEFAULT '',
    "intro"            varchar(500) NOT NULL DEFAULT '',
    "description"      text         NOT NULL DEFAULT '',
    "file_size"        int4         NOT NULL DEFAULT 0,
    "origin_file_name" varchar(255) NOT NULL DEFAULT '',
    "create_time"      int4         NOT NULL DEFAULT 0,
    "update_time"      int4         NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "chat_ai_clawbot_user_skill" ("admin_user_id", "skill_name");
CREATE INDEX ON "chat_ai_clawbot_user_skill" ("admin_user_id");

CREATE TABLE "chat_ai_clawbot_skill"
(
    "id"            serial       NOT NULL primary key,
    "admin_user_id" int4         NOT NULL DEFAULT 0,
    "robot_key"     varchar(15)  NOT NULL DEFAULT '',
    "source_type"   int2         NOT NULL DEFAULT 1,
    "user_skill_id" int4         NOT NULL DEFAULT 0,
    "skill_name"    varchar(50)  NOT NULL DEFAULT '',
    "remark_name"   varchar(50)  NOT NULL DEFAULT '',
    "intro"         varchar(500) NOT NULL DEFAULT '',
    "description"   text         NOT NULL DEFAULT '',
    "file_size"     int4         NOT NULL DEFAULT 0,
    "create_time"   int4         NOT NULL DEFAULT 0,
    "update_time"   int4         NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX ON "chat_ai_clawbot_skill" ("source_type", "robot_key", "skill_name");
CREATE INDEX ON "chat_ai_clawbot_skill" ("robot_key", "source_type");
CREATE INDEX ON "chat_ai_clawbot_skill" ("user_skill_id");
CREATE INDEX ON "chat_ai_clawbot_skill" ("admin_user_id");

COMMENT ON TABLE "chat_ai_clawbot_user_skill" IS 'Clawbot-用户技能表';
COMMENT ON COLUMN "chat_ai_clawbot_user_skill"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_clawbot_user_skill"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_clawbot_user_skill"."skill_name" IS '技能名称(目录名,与SKILL.md同名)';
COMMENT ON COLUMN "chat_ai_clawbot_user_skill"."remark_name" IS '显示名称';
COMMENT ON COLUMN "chat_ai_clawbot_user_skill"."intro" IS '简短介绍';
COMMENT ON COLUMN "chat_ai_clawbot_user_skill"."description" IS '描述(从SKILL.md frontmatter解析)';
COMMENT ON COLUMN "chat_ai_clawbot_user_skill"."file_size" IS '原始zip文件大小(字节)';
COMMENT ON COLUMN "chat_ai_clawbot_user_skill"."origin_file_name" IS '原始上传zip文件名';
COMMENT ON COLUMN "chat_ai_clawbot_user_skill"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_clawbot_user_skill"."update_time" IS '更新时间';

COMMENT ON TABLE "chat_ai_clawbot_skill" IS 'Clawbot-机器人技能表';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."robot_key" IS '机器人key';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."source_type" IS '来源类型:1用户上传,2技能市场';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."user_skill_id" IS '上传技能对应的用户技能ID';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."skill_name" IS '技能名称(目录名,与SKILL.md同名)';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."remark_name" IS '显示名称';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."intro" IS '简短介绍';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."description" IS '描述(从SKILL.md frontmatter解析)';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."file_size" IS '原始zip文件大小(字节)';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_clawbot_skill"."update_time" IS '更新时间';
-- +goose Down
DROP TABLE IF EXISTS "chat_ai_clawbot_skill";
DROP TABLE IF EXISTS "chat_ai_clawbot_user_skill";
