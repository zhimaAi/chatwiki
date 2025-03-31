-- +goose Up

DROP INDEX IF EXISTS chat_ai_library_file_data_bm25;

CREATE INDEX if not exists chat_ai_library_file_data_bm25 ON chat_ai_library_file_data
    using bm25 (id, content, question, answer, library_id)
    WITH (
    key_field = 'id',
    text_fields = '{
        "content": {
          "tokenizer": {"type": "chinese_lindera"}
        },
        "question": {
          "tokenizer": {"type": "chinese_lindera"}
        },
        "answer": {
          "tokenizer": {"type": "chinese_lindera"}
        }
    }'
);