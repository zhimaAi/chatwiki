-- +goose Up

CREATE EXTENSION if not exists age;

SELECT ag_catalog.create_graph('graphrag');