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

func GetLang(db *sql.DB) (string, error) {
    const QUERY = "SELECT lang FROM Config"
    var lang string
    row := db.QueryRow(QUERY)
    if row.Err() != nil {
        return "", row.Err()
    }
    if err := row.Scan(&lang); err != nil {
        return "", err
    }
    return lang, nil
}

func SetLang(db *sql.DB, lang string) error {
    const QUERY = "UPDATE Config SET lang = ?"
    _, err := db.Exec(QUERY, lang)
    return err
}

func AddCSS (db *sql.DB, filename string) error {
    return InsertStringInto(db, "ConfigCSS", "filename", filename)
}

func RemoveCSS(db *sql.DB, filename string) error {
    return RemoveStringFrom(db, "ConfigCSS", "filename", filename)
}

func GetCSS(db *sql.DB) ([]string, error) {
    const QUERY = "SELECT filename FROM ConfigCSS"
    var rows *sql.Rows
    var err error
    rows, err = db.Query(QUERY)
    defer rows.Close()
    if rows, err = db.Query(QUERY); err != nil {
        return []string{}, err
    } else {
        var filename string
        list := make([]string, 0)
        for rows.Next() {
            if errRows := rows.Scan(&filename); errRows != nil {
                return []string{}, errRows
            } else {
                list = append(list, filename)
            }
        }
        return list, nil
    }
    return []string{}, nil
}
