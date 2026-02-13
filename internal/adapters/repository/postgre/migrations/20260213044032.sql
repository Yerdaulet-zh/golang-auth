-- Modify "user_sessions" table
ALTER TABLE "user_sessions" ALTER COLUMN "id" SET DEFAULT gen_random_uuid(), ALTER COLUMN "user_id" SET NOT NULL;
-- Create "audit_user_sessions" table
CREATE TABLE "audit_user_sessions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "session_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "event_type" "audit_event_type" NOT NULL,
  "old_value" text NULL,
  "new_value" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
