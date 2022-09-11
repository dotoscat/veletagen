package manager

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"

)

func OpenDatabase(path string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }
    if _, err := db.Exec(modelDefinition); err != nil {
        return nil, err
    }
    return db, nil
}

func CreateTree(path string) {

}
