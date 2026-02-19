-- Drop index "idx_user_email" from table: "user"
DROP INDEX "idx_user_email";
-- Modify "user" table
ALTER TABLE "user" DROP CONSTRAINT "uni_user_email";
-- Create index "idx_email_active" to table: "user"
CREATE UNIQUE INDEX "idx_email_active" ON "user" ("email") WHERE (deleted_at IS NULL);
