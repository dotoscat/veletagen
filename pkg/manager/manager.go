package manager

// Create database if not exists
// Check database

import (
    "fmt"
    "path/filepath"
    "database/sql"

    _ "embed"
    _ "github.com/mattn/go-sqlite3"
)

//go:embed model-definition.sql
var modelDefinition string

const DB = "index.db"

func GetPathDB(path string) string {
    return filepath.Join(path, DB)
}

func InsertStringInto(db *sql.DB, table, column, value string) error {
    query := fmt.Sprintf("INSERT INTO %v (%v) VALUES (?)", table, column)
    _, err := db.Exec(query, value)
    return err
}

func SelectAllStringsFrom(db *sql.DB, table, column string) ([]string, error) {
    query := fmt.Sprintf("SELECT %v FROM %v", column, table)
    results := make([]string, 0)
    rows, err := db.Query(query)
    var value string
    defer rows.Close()
    if err != nil {
        return []string{}, err
    }
    for rows.Next() {
        if err := rows.Scan(&value); err != nil {
            return []string{}, err
        } else {
            results = append(results, value)
        }
    }
    return results, nil
}

func RemoveStringFrom(db *sql.DB, table, column, value string) error {
    query := fmt.Sprintf("DELETE FROM %v WHERE %v = ?", table, column)
    _, err := db.Exec(query, value)
    return err
}
