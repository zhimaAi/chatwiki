-- +goose Up

CREATE INDEX ON "public"."chat_ai_library_file_data_index" ("delete_time", "status", "library_id");
CREATE INDEX ON "public"."chat_ai_library_file_data_index" ("library_id");
CREATE INDEX ON "public"."chat_ai_library_file_data_index" ("file_id");
CREATE INDEX ON "public"."chat_ai_library_file_data_index" ("status");

CREATE INDEX ON "public"."chat_ai_library_file_data" ("delete_time", "id");