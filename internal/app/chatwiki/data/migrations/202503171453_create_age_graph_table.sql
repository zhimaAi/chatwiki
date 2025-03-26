-- +goose Up

CREATE EXTENSION if not exists age;

SELECT create_graph('graphrag');