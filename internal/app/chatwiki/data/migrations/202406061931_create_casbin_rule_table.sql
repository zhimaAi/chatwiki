-- +goose Up
--
-- CREATE TABLE IF NOT EXISTS "casbin_rule"
-- (
--     id bigserial NOT NULL PRIMARY KEY,
--     ptype varchar(100) NULL DEFAULT '',
--     v0 varchar(100) NULL DEFAULT '',
--     v1 varchar(100) NULL DEFAULT '',
--     v2 varchar(100) NULL DEFAULT '',
--     v3 varchar(100) NULL DEFAULT '',
--     v4 varchar(100) NULL DEFAULT '',
--     v5 varchar(100) NULL DEFAULT ''
-- );
-- CREATE UNIQUE INDEX idx_casbin_rule ON casbin_rule USING btree (ptype, v0, v1, v2, v3, v4, v5);

INSERT INTO public.casbin_rule(ptype, v0, v1, v2, v3, v4, v5) VALUES('g', 'admin', '1', '', '', '', '');
INSERT INTO public.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES('p', '2', 'SystemManage', 'GET', '', '', '');
INSERT INTO public.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES('p', '2', 'RobotManage', 'GET', '', '', '');
INSERT INTO public.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES('p', '2', 'LibraryManage', 'GET', '', '', '');
INSERT INTO public.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES('p', '3', 'RobotManage', 'GET', '', '', '');
INSERT INTO public.casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES('p', '3', 'LibraryManage', 'GET', '', '', '');
