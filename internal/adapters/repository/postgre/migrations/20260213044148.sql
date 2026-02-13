-- Modify "audit_user_sessions" table
ALTER TABLE "audit_user_sessions" ADD CONSTRAINT "fk_audit_user_sessions_user" FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
