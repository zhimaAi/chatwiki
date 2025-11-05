-- +goose Up

ALTER TABLE "public"."chat_ai_library_file"
    ADD COLUMN "father_chunk_paragraph_type" int2         NOT NULL DEFAULT 2,
    ADD COLUMN "father_chunk_separators_no"  varchar(100) NOT NULL DEFAULT '12,11',
    ADD COLUMN "father_chunk_chunk_size"     int4         NOT NULL DEFAULT 1024,
    ADD COLUMN "son_chunk_separators_no"     varchar(100) NOT NULL DEFAULT '8,10',
    ADD COLUMN "son_chunk_chunk_size"        int4         NOT NULL DEFAULT 512;
COMMENT ON COLUMN "public"."chat_ai_library_file"."father_chunk_paragraph_type" IS '父子分段-父块-分段类型:1全文,2段落';
COMMENT ON COLUMN "public"."chat_ai_library_file"."father_chunk_separators_no" IS '父子分段-父块-分段标识符:默认空行,回车';
COMMENT ON COLUMN "public"."chat_ai_library_file"."father_chunk_chunk_size" IS '父子分段-父块-分段最大长度:默认1024';
COMMENT ON COLUMN "public"."chat_ai_library_file"."son_chunk_separators_no" IS '父子分段-子块-分段标识符:默认分号,句号';
COMMENT ON COLUMN "public"."chat_ai_library_file"."father_chunk_chunk_size" IS '父子分段-子块-分段最大长度:默认1024';

ALTER TABLE "public"."chat_ai_library"
    ADD COLUMN "father_chunk_paragraph_type" int2         NOT NULL DEFAULT 2,
    ADD COLUMN "father_chunk_separators_no"  varchar(100) NOT NULL DEFAULT '12,11',
    ADD COLUMN "father_chunk_chunk_size"     int4         NOT NULL DEFAULT 1024,
    ADD COLUMN "son_chunk_separators_no"     varchar(100) NOT NULL DEFAULT '8,10',
    ADD COLUMN "son_chunk_chunk_size"        int4         NOT NULL DEFAULT 512;
COMMENT ON COLUMN "public"."chat_ai_library"."father_chunk_paragraph_type" IS '父子分段-父块-分段类型:1全文,2段落';
COMMENT ON COLUMN "public"."chat_ai_library"."father_chunk_separators_no" IS '父子分段-父块-分段标识符:默认空行,回车';
COMMENT ON COLUMN "public"."chat_ai_library"."father_chunk_chunk_size" IS '父子分段-父块-分段最大长度:默认1024';
COMMENT ON COLUMN "public"."chat_ai_library"."son_chunk_separators_no" IS '父子分段-子块-分段标识符:默认分号,句号';
COMMENT ON COLUMN "public"."chat_ai_library"."father_chunk_chunk_size" IS '父子分段-子块-分段最大长度:默认1024';

ALTER TABLE "public"."chat_ai_library_file_data"
    ADD COLUMN "father_chunk_paragraph_number" int4 NOT NULL DEFAULT 0;
COMMENT ON COLUMN "public"."chat_ai_library_file_data"."father_chunk_paragraph_number" IS '父子分段-父块分段序号';

CREATE INDEX ON "chat_ai_library_file_data" ("father_chunk_paragraph_number", "number");

ALTER TABLE "public"."chat_ai_library"
    ALTER COLUMN "normal_chunk_default_separators_no" SET DEFAULT '12,11';