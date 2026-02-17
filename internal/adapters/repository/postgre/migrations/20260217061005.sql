-- Modify "audit_user_sessions" table
ALTER TABLE "audit_user_sessions" DROP COLUMN "old_value";
-- Modify "user_sessions" table
ALTER TABLE "user_sessions" ADD COLUMN "token" character varying(255) NOT NULL;
