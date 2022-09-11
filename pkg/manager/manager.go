package manager

// Create database if not exists
// Check database

import (
    _ "embed"
)

//go:embed model-definition.sql
var modelDefinition string
