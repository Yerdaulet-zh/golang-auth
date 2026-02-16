-- Create "user_verification" table
CREATE TABLE "user_verification" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "token" character varying(255) NOT NULL,
  "status" "user_verification_status" NOT NULL DEFAULT 'pending',
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_verification_user" FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_user_verification_token" to table: "user_verification"
CREATE UNIQUE INDEX "idx_user_verification_token" ON "user_verification" ("token");
-- Create index "idx_user_verification_user_id" to table: "user_verification"
CREATE INDEX "idx_user_verification_user_id" ON "user_verification" ("user_id");
