package manager

import (
    "database/sql"
    "path/filepath"
    "os"
    "errors"
    "io/fs"

    _ "github.com/mattn/go-sqlite3"
)

func OpenDatabase(path string) (*sql.DB, error) {
    // Check if database exists, if not exists then execute model definition
    var execModelDefinition bool
    _, statErr := os.Stat(path)
    if errors.Is(statErr, fs.ErrNotExist) == true {
        execModelDefinition = true
    } else if statErr != nil {
        return nil, statErr
    }
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }
    if execModelDefinition == true {
        if _, err := db.Exec(modelDefinition); err != nil {
            return nil, err
        }
    }
    return db, nil
}

var CSS_PATH string = filepath.Join("assets", "css")
var SCRIPTS_PATH string = filepath.Join("assets", "scripts")

func CreateTree(path string) error {
    postsPath := filepath.Join(path, "posts")
    scriptsPath := filepath.Join(path, SCRIPTS_PATH)
    cssPath := filepath.Join(path, CSS_PATH)
    if err := os.MkdirAll(postsPath, fs.ModeDir); err != nil {
        return err
    }
    if err := os.MkdirAll(scriptsPath, fs.ModeDir); err != nil {
        return err
    }
    if err := os.MkdirAll(cssPath, fs.ModeDir); err != nil {
        return err
    }
    return nil
}
