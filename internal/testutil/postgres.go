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

// SetupTestDB initializes a Testcontainer Postgres instance, applies Atlas migrations,
// and returns a GORM DB connection ready for repository testing.
func SetupTestDB(t *testing.T) *gorm.DB {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 1. Spin up the Postgres 17 container
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
		t.Fatalf("failed to start container: %s", err)
	}

	// 2. Ensure cleanup
	t.Cleanup(func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	// 3. Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable", "search_path=public")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// 4. Run Atlas Migrations
	applyAtlasMigrations(t, connStr)

	// 5. Initialize GORM with the container's connection string
	// gormDB, err := gorm.Open(gormPostgres.Open(connStr), &gorm.Config{
	// 	// Optional: Disable logging for cleaner test output,
	// 	// unless you need to debug SQL queries.
	// 	// Logger: logger.Default.LogMode(logger.Silent),
	// })
	gormDB, err := gorm.Open(gormPostgres.Open(connStr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Use "user" instead of "users"
		},
	})
	if err != nil {
		t.Fatalf("failed to open gorm.DB: %v", err)
	}

	var tables []string
	gormDB.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables)
	t.Logf("DEBUG: Tables found in DB: %v", tables)

	return gormDB
}

func applyAtlasMigrations(t *testing.T, connStr string) {
	ctx := context.Background()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	migrationDir := filepath.Join(basepath, "..", "adapters", "repository", "postgre", "migrations")

	// Initialize Atlas client
	client, err := atlasexec.NewClient(migrationDir, "atlas")
	if err != nil {
		t.Fatalf("failed to initialize atlas client: %v", err)
	}

	// Use MigrateApply with explicit URL and DirURL.
	// We omit 'Env' to ensure it doesn't try to use your HLC 'local' settings.
	_, err = client.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		URL:    connStr,
		DirURL: "file://.",
	})
	if err != nil {
		t.Fatalf("failed to apply atlas migrations: %v\nCheck if atlas.sum exists in %s", err, migrationDir)
	}

	t.Log("Atlas migrations applied successfully.")
}
