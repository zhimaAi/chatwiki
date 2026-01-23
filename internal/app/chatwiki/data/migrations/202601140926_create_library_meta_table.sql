-- +goose Up

ALTER TABLE "public"."chat_ai_library" ADD COLUMN "show_meta_source" int4 NOT NULL DEFAULT 0;
COMMENT ON COLUMN "public"."chat_ai_library"."show_meta_source" IS '是否显示内置元数据-来源 1显示 0不显示';

ALTER TABLE "public"."chat_ai_library" ADD COLUMN "show_meta_update_time" int4 NOT NULL DEFAULT 0;
COMMENT ON COLUMN "public"."chat_ai_library"."show_meta_update_time" IS '是否显示内置元数据-更新时间 1显示 0不显示';

ALTER TABLE "public"."chat_ai_library" ADD COLUMN "show_meta_create_time" int4 NOT NULL DEFAULT 0;
COMMENT ON COLUMN "public"."chat_ai_library"."show_meta_create_time" IS '是否显示内置元数据-创建时间 1显示 0不显示';

ALTER TABLE "public"."chat_ai_library" ADD COLUMN "show_meta_group" int4 NOT NULL DEFAULT 0;
COMMENT ON COLUMN "public"."chat_ai_library"."show_meta_group" IS '是否显示内置元数据-分组 1显示 0不显示';

create table "public"."library_meta_schema"
(
    "id"                serial NOT NULL PRIMARY KEY,
    "create_time"       int4 NOT NULL DEFAULT 0,
    "update_time"       int4 NOT NULL DEFAULT 0,
    "admin_user_id"     int4 NOT NULL DEFAULT 0,
    "library_id"        int4 NOT NULL DEFAULT 0,
    "name"              varchar(100) NOT NULL DEFAULT 0,
    "key"               varchar(100) NOT NULL,
    "type"              int4 NOT NULL DEFAULT 0,
    "is_show"           int4 NOT NULL DEFAULT 0
);
COMMENT ON TABLE "public"."library_meta_schema" IS '知识库元信息';
COMMENT ON COLUMN "public"."library_meta_schema"."library_id" IS '知识库ID';
COMMENT ON COLUMN "public"."library_meta_schema"."name" IS '数据名';
COMMENT ON COLUMN "public"."library_meta_schema"."key" IS '程序访问的键名';
COMMENT ON COLUMN "public"."library_meta_schema"."type" IS '数据类型';
COMMENT ON COLUMN "public"."library_meta_schema"."is_show" IS '是否显示 1 显示 0 不显示';
CREATE INDEX on "public"."library_meta_schema" USING btree ("admin_user_id", "library_id");
CREATE UNIQUE INDEX ON "public"."library_meta_schema" ("library_id", "key");

ALTER TABLE "public"."chat_ai_library_file" ADD COLUMN "metadata" jsonb NOT NULL DEFAULT '{}'::jsonb;
COMMENT ON COLUMN "public"."chat_ai_library_file"."metadata" IS '元数据';
CREATE INDEX ON "public"."chat_ai_library_file" USING gin ("metadata");

ALTER TABLE "public"."chat_ai_library_file_data" ADD COLUMN "metadata" jsonb NOT NULL DEFAULT '{}'::jsonb;
COMMENT ON COLUMN "public"."chat_ai_library_file_data"."metadata" IS '问答对的元数据';
CREATE INDEX ON "public"."chat_ai_library_file_data" USING gin ("metadata");

ALTER TABLE "public"."chat_ai_robot" ADD COLUMN "meta_search_switch" int4 NOT NULL DEFAULT 0;
COMMENT ON COLUMN "public"."chat_ai_robot"."meta_search_switch" IS '是否启用元信息搜索 1启用 0不启用';

ALTER TABLE "public"."chat_ai_robot" ADD COLUMN "meta_search_type" int2 NOT NULL DEFAULT 1;
COMMENT ON COLUMN "public"."chat_ai_robot"."meta_search_type" IS '元信息搜索条件类型 1且 2或';

ALTER TABLE "public"."chat_ai_robot" ADD COLUMN "meta_search_condition_list" jsonb NOT NULL DEFAULT '{}'::jsonb;;
COMMENT ON COLUMN "public"."chat_ai_robot"."meta_search_condition_list" IS '元信息搜索条件列表';