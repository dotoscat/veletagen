package manager

import (
    "database/sql"

    _ "github.com/mattn/go-sqlite3"
)

func GetTitle(db *sql.DB) (string, error) {
    const QUERY = "SELECT title FROM Config"
    var title string
    row := db.QueryRow(QUERY)
    if row.Err() != nil {
        return "", row.Err()
    }
    if err := row.Scan(&title); err != nil {
        return "", err
    }
    return title, nil
}

func SetTitle(db *sql.DB, title string) error {
    const QUERY = "UPDATE Config SET title = ?"
    _, err := db.Exec(QUERY, title)
    return err
}

func GetPostsPerPage(db *sql.DB) (int64, error) {
    const QUERY = "SELECT posts_per_page FROM Config"
    var postsPerPage int64
    row := db.QueryRow(QUERY)
    if row.Err() != nil {
        return 0, row.Err()
    }
    if err := row.Scan(&postsPerPage); err != nil {
        return 0, err
    }
    return postsPerPage, nil
}

func SetPostsPerPage(db *sql.DB, postsPerPage int64) error {
    const QUERY = "UPDATE Config SET posts_per_page = ?"
    _, err := db.Exec(QUERY, postsPerPage)
    return err
}
