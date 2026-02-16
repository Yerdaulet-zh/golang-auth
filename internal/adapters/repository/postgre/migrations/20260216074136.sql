DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_verification_status') THEN
        CREATE TYPE user_verification_status AS ENUM ('pending', 'consumed', 'invalidated', 'expired');
    END IF;
END $$;
