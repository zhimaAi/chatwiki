-- +goose Up

ALTER TABLE chat_ai_library_file_data
ALTER COLUMN content TYPE VARCHAR(10000),
ALTER COLUMN question TYPE VARCHAR(10000),
ALTER COLUMN answer TYPE VARCHAR(10000);

ALTER TABLE chat_ai_library_file_data_index
ALTER COLUMN content TYPE VARCHAR(10000);