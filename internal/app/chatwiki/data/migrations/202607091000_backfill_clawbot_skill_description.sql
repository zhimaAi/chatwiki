-- +goose Up
-- backfill description from the legacy intro column for rows created before
-- description became the single source of truth; intro is no longer written
UPDATE "chat_ai_clawbot_user_skill" SET "description" = "intro" WHERE "description" = '' AND "intro" <> '';
UPDATE "chat_ai_clawbot_skill" SET "description" = "intro" WHERE "description" = '' AND "intro" <> '';

-- +goose Down
-- not reversible: backfilled values cannot be distinguished from originals
