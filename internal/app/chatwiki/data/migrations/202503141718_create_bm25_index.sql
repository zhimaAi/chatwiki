-- +goose Up

CREATE EXTENSION if not exists pg_search;

CREATE INDEX if not exists chat_ai_library_file_data_bm25 ON chat_ai_library_file_data
using bm25 (id, content, question, answer, library_id)
WITH (
    key_field = 'id',
    text_fields = '{
        "content": {
          "tokenizer": {"type": "chinese_lindera"}
        }
    }'
);