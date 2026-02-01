CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin','manager','viewer'))
);

INSERT INTO users(username, password_hash, role) VALUES
('admin', '$2a$12$examplehashadmin', 'admin'),
('manager', '$2a$12$examplehashmanager', 'manager'),
('viewer', '$2a$12$examplehashviewer', 'viewer');


CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    quantity INT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS item_history (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,       
    action TEXT NOT NULL,       
    changed_by TEXT NOT NULL,   
    old_data JSONB,
    new_data JSONB,
    changed_at TIMESTAMP DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION trg_log_item_history()
RETURNS TRIGGER AS $$
DECLARE
    username TEXT;
BEGIN
    username := current_setting('myapp.current_user', true);
    IF username IS NULL THEN
        username := 'system';
    END IF;

    IF TG_OP = 'INSERT' THEN
        INSERT INTO item_history(item_id, action, changed_by, old_data, new_data)
        VALUES (NEW.id, 'INSERT', username, NULL, to_jsonb(NEW));
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO item_history(item_id, action, changed_by, old_data, new_data)
        VALUES (NEW.id, 'UPDATE', username, to_jsonb(OLD), to_jsonb(NEW));
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO item_history(item_id, action, changed_by, old_data, new_data)
        VALUES (OLD.id, 'DELETE', username, to_jsonb(OLD), NULL);
        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER trg_items_history
BEFORE INSERT OR UPDATE OR DELETE
ON items
FOR EACH ROW
EXECUTE FUNCTION trg_log_item_history();
