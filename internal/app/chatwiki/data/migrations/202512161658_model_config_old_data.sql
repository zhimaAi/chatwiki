-- +goose Up

UPDATE "public"."chat_ai_model_config"
SET model_types = 'LLM,TEXT EMBEDDING'
WHERE model_define = 'openaiAgent';

UPDATE "public"."chat_ai_model_config"
SET model_types = 'LLM,TEXT EMBEDDING'
WHERE model_define = 'azure';

UPDATE "public"."chat_ai_model_config"
SET model_types = 'LLM,TEXT EMBEDDING'
WHERE model_define = 'ollama';

UPDATE "public"."chat_ai_model_config"
SET model_types = 'LLM,TEXT EMBEDDING,RERANK'
WHERE model_define = 'xinference';

UPDATE "public"."chat_ai_model_config"
SET model_types = 'LLM,TEXT EMBEDDING,IMAGE'
WHERE model_define = 'doubao';

UPDATE "public"."chat_ai_model_config"
SET model_types = 'LLM,TTS'
WHERE model_define = 'minimax';
