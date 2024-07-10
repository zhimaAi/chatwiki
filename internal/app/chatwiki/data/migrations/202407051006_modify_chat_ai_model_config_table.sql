-- +goose Up

ALTER TABLE chat_ai_model_config
    ALTER COLUMN app_id TYPE VARCHAR(32) USING app_id::text,
    ALTER COLUMN app_id SET DEFAULT '';
