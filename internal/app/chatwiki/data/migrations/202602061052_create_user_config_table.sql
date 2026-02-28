-- +goose Up

CREATE TABLE "user_config"
(
    "id"                    serial       NOT NULL primary key,
    "admin_user_id"         int4         NOT NULL DEFAULT 0,
    "user_id"               int4         NOT NULL DEFAULT 0,
    "qa_merge_similarity"   numeric(5,4) NOT NULL DEFAULT 0.8,
    "create_time"           int4         NOT NULL DEFAULT 0,
    "update_time"           int4         NOT NULL DEFAULT 0
);

CREATE INDEX ON "user_config" ("admin_user_id");
CREATE INDEX ON "user_config" ("user_id");

COMMENT ON TABLE "user_config" IS 'User configuration table';
COMMENT ON COLUMN "user_config"."id" IS 'ID';
COMMENT ON COLUMN "user_config"."admin_user_id" IS 'Admin user ID';
COMMENT ON COLUMN "user_config"."user_id" IS 'User ID';
COMMENT ON COLUMN "user_config"."qa_merge_similarity" IS 'QA merge similarity threshold (0-1)';
COMMENT ON COLUMN "user_config"."create_time" IS 'Create time';
COMMENT ON COLUMN "user_config"."update_time" IS 'Update time';

-- +goose Down
DROP TABLE "user_config";
