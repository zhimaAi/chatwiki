-- +goose Up

CREATE EXTENSION IF NOT EXISTS pg_search;

-- +goose StatementBegin
DO $$
DECLARE
    tokenizer regtype;
    index_method text := 'bm25';
    installed_version text;
    package_version text;
    upgrade_path text;
    target_version constant text := '0.22.6';
BEGIN
    SELECT extversion INTO installed_version FROM pg_extension WHERE extname = 'pg_search';
    IF installed_version IS NULL THEN
        RAISE EXCEPTION 'pg_search must be installed and enabled before creating the library BM25 index';
    END IF;

    -- Existing data volumes retain their extension version after the Docker image is rebuilt.
    -- Use the package's official update scripts; never drop or recreate the extension.
    IF installed_version <> target_version THEN
        IF installed_version !~ '^[0-9]+\.[0-9]+\.[0-9]+$' THEN
            RAISE EXCEPTION 'Unsupported pg_search version format: %', installed_version;
        END IF;
        IF string_to_array(installed_version, '.')::integer[] < ARRAY[0, 22, 6] THEN
            SELECT path INTO upgrade_path
            FROM pg_extension_update_paths('pg_search')
            WHERE source = installed_version AND target = target_version;
            IF upgrade_path IS NULL THEN
                RAISE EXCEPTION 'No pg_search upgrade path from % to %. Rebuild the PostgreSQL Docker image with pg_search %',
                    installed_version, target_version, target_version;
            END IF;
            EXECUTE format('ALTER EXTENSION pg_search UPDATE TO %L', target_version);
        ELSE
            -- A newer database must not be used with an older extension binary.
            SELECT default_version INTO package_version
            FROM pg_available_extensions WHERE name = 'pg_search';
            IF package_version ~ '^[0-9]+\.[0-9]+\.[0-9]+$' THEN
                IF string_to_array(package_version, '.')::integer[] < string_to_array(installed_version, '.')::integer[] THEN
                    RAISE EXCEPTION 'Installed pg_search % is newer than Docker package %. Use a matching or newer PostgreSQL image',
                        installed_version, package_version;
                END IF;
            END IF;
        END IF;
    END IF;

    tokenizer := to_regtype('pdb.jieba_yada');
    IF tokenizer IS NULL THEN
        tokenizer := to_regtype('pdb.jieba');
    END IF;
    IF tokenizer IS NULL THEN
        RAISE EXCEPTION 'pg_search requires pdb.jieba_yada or pdb.jieba for library BM25 search';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_am WHERE amname = 'bm25') THEN
        RAISE EXCEPTION 'pg_search requires the bm25 index access method';
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_proc AS p
        JOIN pg_namespace AS n ON n.oid = p.pronamespace
        WHERE n.nspname = 'pdb' AND p.proname = 'score' AND p.pronargs = 1
    ) OR NOT EXISTS (SELECT 1 FROM pg_operator WHERE oprname = '|||') THEN
        RAISE EXCEPTION 'pg_search requires pdb.score(key_field) and the ||| match operator';
    END IF;

    EXECUTE format(
        'CREATE INDEX chat_ai_library_file_data_index_bm25
         ON public.chat_ai_library_file_data_index
         USING %I (id, (content::%s), library_id, data_id, delete_time, file_id, admin_user_id, type)
         WITH (key_field = ''id'')',
        index_method, tokenizer
    );
END;
$$;
-- +goose StatementEnd

-- +goose Down
DROP INDEX IF EXISTS public.chat_ai_library_file_data_index_bm25;
