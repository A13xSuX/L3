-- +goose Up

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION audit_items_changes()
    RETURNS TRIGGER AS $$
DECLARE
    v_user_id_text TEXT;
    v_username TEXT;
    v_role TEXT;
BEGIN
    v_user_id_text := current_setting('app.current_user_id', true);
    v_username := current_setting('app.current_username', true);
    v_role := current_setting('app.current_role', true);

    IF TG_OP = 'INSERT' THEN
        INSERT INTO audit (
            item_id,
            action,
            old_data,
            new_data,
            changed_by_user_id,
            changed_by_username,
            changed_by_role,
            changed_at
        )
        VALUES (
                   NEW.id,
                   'INSERT',
                   NULL,
                   to_jsonb(NEW),
                   CASE
                       WHEN v_user_id_text IS NULL OR v_user_id_text = '' THEN NULL
                       ELSE v_user_id_text::BIGINT
                       END,
                   COALESCE(v_username, 'unknown'),
                   COALESCE(v_role, 'viewer'),
                   NOW()
               );

        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO audit (
            item_id,
            action,
            old_data,
            new_data,
            changed_by_user_id,
            changed_by_username,
            changed_by_role,
            changed_at
        )
        VALUES (
                   NEW.id,
                   'UPDATE',
                   to_jsonb(OLD),
                   to_jsonb(NEW),
                   CASE
                       WHEN v_user_id_text IS NULL OR v_user_id_text = '' THEN NULL
                       ELSE v_user_id_text::BIGINT
                       END,
                   COALESCE(v_username, 'unknown'),
                   COALESCE(v_role, 'viewer'),
                   NOW()
               );

        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO audit (
            item_id,
            action,
            old_data,
            new_data,
            changed_by_user_id,
            changed_by_username,
            changed_by_role,
            changed_at
        )
        VALUES (
                   OLD.id,
                   'DELETE',
                   to_jsonb(OLD),
                   NULL,
                   CASE
                       WHEN v_user_id_text IS NULL OR v_user_id_text = '' THEN NULL
                       ELSE v_user_id_text::BIGINT
                       END,
                   COALESCE(v_username, 'unknown'),
                   COALESCE(v_role, 'viewer'),
                   NOW()
               );

        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS trg_audit_items_changes ON items;

CREATE TRIGGER trg_audit_items_changes
    AFTER INSERT OR UPDATE OR DELETE ON items
    FOR EACH ROW
EXECUTE FUNCTION audit_items_changes();

-- +goose Down

DROP TRIGGER IF EXISTS trg_audit_items_changes ON items;

-- +goose StatementBegin
DROP FUNCTION IF EXISTS audit_items_changes();
-- +goose StatementEnd