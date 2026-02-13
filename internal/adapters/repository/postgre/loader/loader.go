package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	persistency "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
	usersessions "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user_sessions"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	stmts, err := gormschema.New("postgres", gormschema.WithConfig(config)).Load(
		&persistency.User{},
		&usersessions.UserSessions{},
		&usersessions.AuditUserSessions{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
