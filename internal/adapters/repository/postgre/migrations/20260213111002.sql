-- Create "user_credentials" table
CREATE TABLE "user_credentials" (
  "user_id" uuid NOT NULL,
  "password_hash" character varying(255) NOT NULL,
  "last_password_change_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("user_id"),
  CONSTRAINT "fk_user_credentials_user" FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
