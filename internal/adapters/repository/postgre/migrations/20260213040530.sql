-- Create "user" table
CREATE TABLE "user" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "email" character varying(255) NOT NULL,
  "user_status" "user_status" NOT NULL DEFAULT 'pending_verification',
  "is_mfa_enabled" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_user_email" UNIQUE ("email")
);
-- Create index "idx_user_deleted_at" to table: "user"
CREATE INDEX "idx_user_deleted_at" ON "user" ("deleted_at");
-- Create index "idx_user_email" to table: "user"
CREATE INDEX "idx_user_email" ON "user" ("email");
