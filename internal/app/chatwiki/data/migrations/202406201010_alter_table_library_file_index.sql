-- +goose Up

alter table chat_ai_library_file_data_index add column embedding_var  vector;
update chat_ai_library_file_data_index set embedding_var = embedding  where id>0;
alter table chat_ai_library_file_data_index drop column embedding;
alter table chat_ai_library_file_data_index rename column embedding_var to embedding;
COMMENT ON COLUMN "chat_ai_library_file_data_index"."embedding" IS '转换为可变维度向量的文档';
