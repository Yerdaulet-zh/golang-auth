-- Create "user_sessions" table
CREATE TABLE "user_sessions" (
  "id" uuid NOT NULL,
  "user_id" uuid NULL,
  "ip_address" inet NOT NULL,
  "user_agent" text NOT NULL,
  "device" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "last_active" timestamptz NOT NULL DEFAULT now(),
  "expires_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_sessions_user" FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
