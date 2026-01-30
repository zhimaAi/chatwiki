-- +goose Up

UPDATE public.chat_ai_model_config
SET model_types = 'LLM,TEXT EMBEDDING,RERANK,IMAGE'
WHERE model_define = 'tongyi';
