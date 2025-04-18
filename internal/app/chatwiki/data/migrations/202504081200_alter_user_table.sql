-- +goose Up
COMMENT ON COLUMN "public"."chat_ai_library_file"."status" IS '状态:0待转换PDF,4待用户切分,5待爬取,6,爬取中,7,爬取失败,1学习中,2学习完成,3文件异常,8:部分转换异常';