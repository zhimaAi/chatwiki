-- +goose Up

ALTER TABLE "public"."chat_ai_robot"
    ADD COLUMN "show_ai_msg_gzh"  int2 NOT NULL DEFAULT 0,
    ADD COLUMN "show_typing_gzh"  int2 NOT NULL DEFAULT 0,
    ADD COLUMN "show_ai_msg_mini" int2 NOT NULL DEFAULT 0,
    ADD COLUMN "show_typing_mini" int2 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."chat_ai_robot"."show_ai_msg_gzh" IS '公众号回复显示内容由AI生成开关:0关,1开';
COMMENT ON COLUMN "public"."chat_ai_robot"."show_typing_gzh" IS '公众号回复时显示正在输入中开关:0关,1开';
COMMENT ON COLUMN "public"."chat_ai_robot"."show_ai_msg_mini" IS '小程序回复显示内容由AI生成开关:0关,1开';
COMMENT ON COLUMN "public"."chat_ai_robot"."show_typing_mini" IS '小程序回复时显示正在输入中开关:0关,1开';
