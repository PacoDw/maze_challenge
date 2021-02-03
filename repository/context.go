package repository

import (
	"context"

	"github.com/spf13/cast"
)

var (
	// ContextDBName represent the value that contains the current context.
	ContextDBName = contextKey("DB_NAME")
)

type contextKey string

// DBNameSet can be used to set the db name to the current context.
func DBNameSet(ctx context.Context, dbName string) context.Context {
	return context.WithValue(ctx, string(ContextDBName), dbName)
}

// DBName retrieves the db name that exists in current context.
func DBName(ctx context.Context) string {
	return cast.ToString(ctx.Value(string(ContextDBName)))
}
