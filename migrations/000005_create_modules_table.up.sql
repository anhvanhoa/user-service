CREATE TYPE module_status AS ENUM ('active', 'inactive');

CREATE TABLE modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status module_status DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_modules_updated_at BEFORE
UPDATE ON modules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column ();