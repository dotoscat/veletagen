//   Copyright 2023 Oscar Triano García
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.package manager

// Create database if not exists
// Check database
package manager

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

func DoesStringExistIn(db *sql.DB, table, column, value string) (bool, error) {
    query := fmt.Sprintf("SELECT COUNT(*) > 0 FROM %v WHERE %v = ?", table, column)
    row := db.QueryRow(query, value)
    var exists bool
    err := row.Scan(&exists)
    return exists, err
}

func GetCategories(db *sql.DB) ([]string, error) {
    const QUERY = "SELECT name FROM Category"
    categories := make([]string, 0)
    rows, err := db.Query(QUERY)
    defer rows.Close()
    var name string
    if err != nil {
        return categories, nil
    }
    for rows.Next() {
        if err := rows.Scan(&name); err != nil {
            return categories, err
        } else {
            categories = append(categories, name)
        }
    }
    return categories, nil
}

