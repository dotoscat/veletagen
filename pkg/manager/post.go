package manager

import (
    "database/sql"
)

func AddPost(db *sql.DB, filename string) error {
    return InsertStringInto(db, "Post", "filename", filename)
}

func RemovePost(db *sql.DB, filename string) error {
    return RemoveStringFrom(db, "Post", "filename", filename)
}
