package manager

// Create database if not exists
// Check database

import (
    "path/filepath"

    _ "embed"
)

//go:embed model-definition.sql
var modelDefinition string

const DB = "index.db"

func GetPathDB(path string) string {
    return filepath.Join(path, DB)
}
