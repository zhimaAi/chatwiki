-- +goose Up

CREATE INDEX on "public"."chat_ai_stat_library_robot_tip"("library_id","admin_user_id");

-- clean stat library data robot tip
DELETE FROM "public"."chat_ai_stat_library_data_robot_tip" T
WHERE NOT EXISTS (
    SELECT 1 FROM "public"."chat_ai_library_file_data" d WHERE d.ID = T.data_id
);

-- clean stat library robot tip
DELETE FROM "public"."chat_ai_stat_library_robot_tip" T
WHERE NOT EXISTS (
    SELECT 1 FROM "public"."chat_ai_library" l WHERE l.ID = T.library_id
);

