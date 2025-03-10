-- +goose Up

-- 启用向量拓展
CREATE EXTENSION IF NOT EXISTS vector;

-- 开启中文分词器扩展
CREATE EXTENSION IF NOT EXISTS zhparser;

-- 配置中文分词器
CREATE TEXT SEARCH CONFIGURATION zhima_zh_parser (PARSER = zhparser);
ALTER TEXT SEARCH CONFIGURATION zhima_zh_parser ADD MAPPING FOR n,v,a,i,e,l WITH simple;



