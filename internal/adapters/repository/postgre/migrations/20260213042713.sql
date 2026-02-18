-- Audit Event Type Enum
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'audit_event_type') THEN
        CREATE TYPE audit_event_type AS ENUM ('IP_CHANGE', 'UA_CHANGE', 'LOGIN', 'LOGOUT', 'CONCURRENCY_LIMIT_REACHED');
    END IF;
END $$;