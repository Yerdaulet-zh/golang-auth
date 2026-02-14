package testutil

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"ariga.io/atlas-go-sdk/atlasexec"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// SetupGlobalTestDB starts a container without a testing.T context.
// Returns the GORM DB, the container instance (for cleanup), and error.
func SetupGlobalTestDB() (*gorm.DB, testcontainers.Container, error) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("enterprise_test_db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable", "search_path=public")
	if err != nil {
		return nil, nil, err
	}

	// Reuse your existing migration logic (wrap it in a check to handle nil *testing.T)
	applyAtlasMigrations(nil, connStr)

	gormDB, err := gorm.Open(gormPostgres.Open(connStr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	return gormDB, pgContainer, err
}

func applyAtlasMigrations(t *testing.T, connStr string) {
	ctx := context.Background()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	migrationDir := filepath.Join(basepath, "..", "adapters", "repository", "postgre", "migrations")

	// Initialize Atlas client
	client, err := atlasexec.NewClient(migrationDir, "atlas")
	if err != nil {
		// Fix 1: Nil check before Fatalf
		if t != nil {
			t.Fatalf("failed to initialize atlas client: %v", err)
		}
		panic("failed to initialize atlas client: " + err.Error())
	}

	_, err = client.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		URL:    connStr,
		DirURL: "file://.",
	})

	if err != nil {
		// Fix 2: Nil check before Fatalf
		if t != nil {
			t.Fatalf("fail: %v", err)
		}
		panic("failed to apply migrations: " + err.Error())
	}

	// Fix 3: Nil check before Log
	if t != nil {
		t.Log("Atlas migrations applied successfully.")
	}
}

// TruncateAllTables wipes all data from the public schema without dropping tables.
func TruncateAllTables(db *gorm.DB) error {
	var tableNames []string

	// Query to get all user-defined tables in the public schema
	err := db.Raw(`
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_type = 'BASE TABLE'
        AND table_name <> 'atlas_schema_revisions'
    `).Scan(&tableNames).Error

	if err != nil {
		return err
	}

	if len(tableNames) == 0 {
		return nil
	}

	// Join table names and execute TRUNCATE CASCADE
	// Using CASCADE ensures foreign key constraints don't block the truncation
	query := "TRUNCATE TABLE "
	for i, name := range tableNames {
		query += "\"" + name + "\""
		if i < len(tableNames)-1 {
			query += ", "
		}
	}
	query += " CASCADE"

	return db.Exec(query).Error
}
