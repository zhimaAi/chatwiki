-- +goose Up


ALTER TABLE "public"."chat_ai_robot" 
  ADD COLUMN "tips_before_answer_switch" bool NOT NULL DEFAULT true;

COMMENT ON COLUMN "public"."chat_ai_robot"."tips_before_answer_switch" IS '答案生成中提示语开关  false关， true开 默认';


ALTER TABLE "public"."chat_ai_robot" 
  ADD COLUMN "tips_before_answer_content" varchar(50) NOT NULL DEFAULT '思考中、请稍等';

COMMENT ON COLUMN "public"."chat_ai_robot"."tips_before_answer_content" IS '答案生成中提示语内容';


