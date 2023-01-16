package manager

// Create database if not exists
// Check database

import (
    "fmt"
    "path/filepath"
    "database/sql"

    _ "embed"
    _ "github.com/mattn/go-sqlite3"

    "github.com/dotoscat/veletagen/pkg/common"
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

func DoesStringExistIn(db *sql.DB, table, column, value string) (bool, error) {
    query := fmt.Sprintf("SELECT COUNT(*) > 0 FROM %v WHERE %v = ?", table, column)
    row := db.QueryRow(query, value)
    var exists bool
    err := row.Scan(&exists)
    return exists, err
}

//func AddManyStringsToOneString(db *sql.DB, table, single string, members []string) error {
//
//}

func GetWebsiteBase(db *sql.DB) (common.WebsiteBase, error) {
    const QUERY = "SELECT title, posts_per_page, output_path, lang, license FROM Config"
    var website common.WebsiteBase
    row := db.QueryRow(QUERY)
    err := row.Scan(
        &website.Title,
        &website.PostsPerPage,
        &website.OutputPath,
        &website.Lang,
        &website.License,
    )
    return website, err
}
