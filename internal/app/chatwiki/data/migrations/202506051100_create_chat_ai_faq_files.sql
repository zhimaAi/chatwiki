-- +goose Up

CREATE TABLE "public"."chat_ai_faq_files" (
    "id" serial NOT NULL primary key,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "creator" int4 NOT NULL DEFAULT 0,
    "file_url" varchar(500) NOT NULL DEFAULT '',
    "file_name" varchar(100) NOT NULL DEFAULT '',
    "status" int2 NOT NULL DEFAULT 0,
    "errmsg" varchar(1000) NOT NULL DEFAULT '',
    "file_ext" varchar(4) NOT NULL DEFAULT '',
    "file_size" int4 NOT NULL DEFAULT 0,
    "word_total" int4 NOT NULL DEFAULT 0,
    "split_total" int4 NOT NULL DEFAULT 0,
    "is_table_file" int2 NOT NULL DEFAULT 0,
    "separators_no" varchar(100) NOT NULL DEFAULT '',
    "chunk_type" int4 NOT NULL DEFAULT 0,
    "chunk_size" int4 NOT NULL DEFAULT 0,
    "chunk_overlap" int4 NOT NULL DEFAULT 0,
    "chunk_model_config_id" int4 NOT NULL DEFAULT 0,
    "chunk_model" varchar(100) NOT NULL DEFAULT '',
    "chunk_prompt" varchar(1000) NOT NULL DEFAULT '',
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_faq_files" ("admin_user_id");

COMMENT ON TABLE "public"."chat_ai_faq_files" IS 'FAQ文件';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."id" IS 'ID';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."admin_user_id" IS '管理员用户ID';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."file_url" IS '文件链接';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."file_name" IS '文件名称';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."status" IS '状态:0:排队中 1:解析中 2:提取中 3:提取完成 4:提取失败';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."errmsg" IS '切分报错信息';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."file_ext" IS '文件格式';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."file_size" IS '文件大小';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."word_total" IS '单词数量';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."split_total" IS '切分数量';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."chunk_type" IS '分段类型 1:长度 2:分隔符';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."separators_no" IS '自定义分段-分隔符序号集';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."chunk_size" IS '自定义分段-分段最大长度';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."create_time" IS '创建时间';

COMMENT ON COLUMN "public"."chat_ai_faq_files"."update_time" IS '更新时间';

CREATE TABLE "public"."chat_ai_faq_files_data" (
    "id" serial NOT NULL primary key,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "file_id" int4 NOT NULL DEFAULT 0,
    "number" int4 NOT NULL DEFAULT 0,
    "page_num" int4 NOT NULL DEFAULT 0,
    "type" int4 NOT NULL DEFAULT 0,
    "title" varchar(100) NOT NULL DEFAULT '',
    "content" varchar(10000) NOT NULL DEFAULT '',
    "err_content" varchar(10000) NOT NULL DEFAULT '',
    "images" json NOT NULL DEFAULT '[]',
    "word_total" int4 NOT NULL DEFAULT 0,
    "split_status" int2 NOT NULL DEFAULT 0,
    "split_errmsg" varchar(1000) NOT NULL DEFAULT '',
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_faq_files_data" ("admin_user_id");

CREATE INDEX ON "public"."chat_ai_faq_files_data" ("file_id");

COMMENT ON TABLE "public"."chat_ai_faq_files_data" IS 'FAQ文件分段表';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."id" IS 'ID';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."admin_user_id" IS '管理员用户ID';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."file_id" IS '文件ID';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."number" IS '文档编号';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."page_num" IS '文档页码';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."type" IS '数据类型 1普通段落 2文档问答 3excel问答';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."title" IS '分段标题';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."word_total" IS '单词数量';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."content" IS '文档段落';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."split_status" IS '0:待切分 1:切分成功 2:切分失败';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."create_time" IS '创建时间';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data"."update_time" IS '更新时间';

CREATE TABLE "public"."chat_ai_faq_files_data_qa" (
    "id" serial NOT NULL primary key,
    "admin_user_id" int4 NOT NULL DEFAULT 0,
    "data_id" int4 NOT NULL DEFAULT 0,
    "file_id" int4 NOT NULL DEFAULT 0,
    "number" int4 NOT NULL DEFAULT 0,
    "page_num" int4 NOT NULL DEFAULT 0,
    "question" varchar(1000) NOT NULL DEFAULT '',
    "answer" varchar(1000) NOT NULL DEFAULT '',
    "images" json NOT NULL DEFAULT '[]',
    "is_import" int4 NOT NULL DEFAULT 0,
    "library_id" int4 NOT NULL DEFAULT 0,
    "create_time" int4 NOT NULL DEFAULT 0,
    "update_time" int4 NOT NULL DEFAULT 0
);

CREATE INDEX ON "public"."chat_ai_faq_files_data_qa" ("admin_user_id");

CREATE INDEX ON "public"."chat_ai_faq_files_data_qa" ("file_id");

COMMENT ON TABLE "public"."chat_ai_faq_files_data_qa" IS 'FAQ文件提取QA表';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data_qa"."is_import" IS '1:已导入';

COMMENT ON COLUMN "public"."chat_ai_faq_files_data_qa"."library_id" IS '导入的知识库id';