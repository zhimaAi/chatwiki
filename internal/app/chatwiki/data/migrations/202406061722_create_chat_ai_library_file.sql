-- +goose Up

CREATE TABLE "chat_ai_library_file"
(
    "id"                serial          NOT NULL primary key,
    "admin_user_id"     int4            NOT NULL DEFAULT 0,
    "library_id"        int4            NOT NULL DEFAULT 0,
    "file_url"          varchar(500)    NOT NULL DEFAULT '',
    "pdf_url"           varchar(500)    NOT NULL DEFAULT '',
    "file_name"         varchar(100)    NOT NULL DEFAULT '',
    "status"            int2            NOT NULL DEFAULT 0,
    "errmsg"            varchar(1000)   NOT NULL DEFAULT '',
    "file_ext"          varchar(4)      NOT NULL DEFAULT '',
    "file_size"         int4            NOT NULL DEFAULT 0,
    "word_total"        int4            NOT NULL DEFAULT 0,
    "split_total"       int4            NOT NULL DEFAULT 0,
    "is_table_file"     int2            NOT NULL DEFAULT 0,
    "is_diy_split"      int2            NOT NULL DEFAULT 0,
    "separators_no"     varchar(100)    NOT NULL DEFAULT '',
    "chunk_size"        int4            NOT NULL DEFAULT 0,
    "chunk_overlap"     int4            NOT NULL DEFAULT 0,
    "is_qa_doc"         int2            NOT NULL DEFAULT 0,
    "question_lable"    varchar(100)    NOT NULL DEFAULT '',
    "answer_lable"      varchar(100)    NOT NULL DEFAULT '',
    "question_column"   varchar(16)     NOT NULL DEFAULT '',
    "answer_column"     varchar(16)     NOT NULL DEFAULT '',
    "qa_index_type"     int4            NOT NULL DEFAULT 1,
    "create_time"       int4            NOT NULL DEFAULT 0,
    "update_time"       int4            NOT NULL DEFAULT 0
);

CREATE INDEX ON "chat_ai_library_file" ("admin_user_id");
CREATE INDEX ON "chat_ai_library_file" ("library_id");

COMMENT ON TABLE "chat_ai_library_file" IS '文档问答机器人-知识库文件';

COMMENT ON COLUMN "chat_ai_library_file"."id" IS 'ID';
COMMENT ON COLUMN "chat_ai_library_file"."admin_user_id" IS '管理员用户ID';
COMMENT ON COLUMN "chat_ai_library_file"."library_id" IS '知识库ID';
COMMENT ON COLUMN "chat_ai_library_file"."file_url" IS '文件链接';
COMMENT ON COLUMN "chat_ai_library_file"."pdf_url" IS 'pdf链接';
COMMENT ON COLUMN "chat_ai_library_file"."file_name" IS '文件名称';
COMMENT ON COLUMN "chat_ai_library_file"."status" IS '状态:0待转换PDF,4待用户切分,1学习中,2学习完成,3文件异常';
COMMENT ON COLUMN "chat_ai_library_file"."errmsg" IS '切分报错信息';
COMMENT ON COLUMN "chat_ai_library_file"."file_ext" IS '文件格式';
COMMENT ON COLUMN "chat_ai_library_file"."file_size" IS '文件大小';
COMMENT ON COLUMN "chat_ai_library_file"."word_total" IS '单词数量';
COMMENT ON COLUMN "chat_ai_library_file"."split_total" IS '切分数量';
COMMENT ON COLUMN "chat_ai_library_file"."is_table_file" IS '是否是表格文件(xlsx,xls,csv)';
COMMENT ON COLUMN "chat_ai_library_file"."is_diy_split" IS '是否自定义分段';
COMMENT ON COLUMN "chat_ai_library_file"."separators_no" IS '自定义分段-分隔符序号集';
COMMENT ON COLUMN "chat_ai_library_file"."chunk_size" IS '自定义分段-分段最大长度';
COMMENT ON COLUMN "chat_ai_library_file"."chunk_overlap" IS '自定义分段-分段重叠长度';
COMMENT ON COLUMN "chat_ai_library_file"."is_qa_doc" IS '是否是QA文档';
COMMENT ON COLUMN "chat_ai_library_file"."question_lable" IS 'QA文档-问题开始标识符';
COMMENT ON COLUMN "chat_ai_library_file"."answer_lable" IS 'QA文档-答案开始标识符';
COMMENT ON COLUMN "chat_ai_library_file"."question_column" IS 'QA文档-问题所在列';
COMMENT ON COLUMN "chat_ai_library_file"."answer_column" IS 'QA文档-答案所在列';
COMMENT ON COLUMN "chat_ai_library_file"."qa_index_type" IS 'QA文档-索引方式,1问题答案一起生成索引,2仅对问题生成索引';
COMMENT ON COLUMN "chat_ai_library_file"."create_time" IS '创建时间';
COMMENT ON COLUMN "chat_ai_library_file"."update_time" IS '更新时间';
