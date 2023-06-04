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
package manager

import (
    "database/sql"

    _ "github.com/mattn/go-sqlite3"
)

type Config struct {
    Title string
    License string
    Lang string
    OutputPath string
    PostsPerPage int64
}

func GetConfig(db *sql.DB) (Config, error) {
    const QUERY = "SELECT title, posts_per_page, output_path, lang, license FROM Config"
    var config Config
    row := db.QueryRow(QUERY)
    err := row.Scan(
        &config.Title,
        &config.PostsPerPage,
        &config.OutputPath,
        &config.Lang,
        &config.License,
    )
    return config, err
}

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

func GetLicense(db *sql.DB) (string, error) {
    const QUERY = "SELECT license FROM Config"
    var license string
    row := db.QueryRow(QUERY)
    if row.Err() != nil {
        return "", row.Err()
    }
    if err := row.Scan(&license); err != nil {
        return "", err
    }
    return license, nil
}

func SetLicense(db *sql.DB, license string) error {
    const QUERY = "UPDATE Config SET license = ?"
    _, err := db.Exec(QUERY, license)
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

func AddScript(db *sql.DB, filename string) error {
    return InsertStringInto(db, "ConfigScript", "filename", filename)
}

func RemoveScript(db *sql.DB, filename string) error {
    return RemoveStringFrom(db, "ConfigScript", "filename", filename)
}

func GetScripts(db *sql.DB) ([]string, error) {
    return SelectAllStringsFrom(db, "ConfigScript", "filename")
}
