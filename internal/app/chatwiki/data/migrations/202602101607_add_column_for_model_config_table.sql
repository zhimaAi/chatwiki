-- +goose Up

create table "public"."model_define_weight" (
    "id" serial not null primary key,
    "admin_user_id" int4 not null default 0,
    "model_config_id" varchar(100) not null default '',
    "weight" int4 not null default 0,
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0
);

-- 添加索引 node_key
CREATE INDEX on "model_define_weight"("admin_user_id" , "model_config_id");

COMMENT ON TABLE "public"."model_define_weight" IS '模型排序';
COMMENT ON COLUMN "public"."model_define_weight"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "public"."model_define_weight"."model_config_id" IS '配置id';
COMMENT ON COLUMN "public"."model_define_weight"."weight" IS '权重';
