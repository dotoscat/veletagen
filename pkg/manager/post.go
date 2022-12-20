package manager

import (
    "database/sql"
    "errors"
    "fmt"
    "log"

    "github.com/dotoscat/veletagen/pkg/common"
)

func AddPost(db *sql.DB, filename string) error {
    return InsertStringInto(db, "Post", "filename", filename)
}

func RemovePost(db *sql.DB, filename string) error {
    return RemoveStringFrom(db, "Post", "filename", filename)
}

func AddTag(db *sql.DB, name string) error {
    if exists, err := DoesStringExistIn(db, "Tag", "name", name); err != nil {
        return err
    } else if exists == true {
        return nil
    }
    return InsertStringInto(db, "Tag", "name", name)
}

func RemoveTag(db *sql.DB, name string) error {
    return RemoveStringFrom(db, "Tag", "name", name)
}

func AddTagsToPost(db *sql.DB, filename string, tags common.Tags) error{
    if exists, err := DoesStringExistIn(db, "Post", "filename", filename); err != nil {
        return err
    } else if exists == false {
        return errors.New(fmt.Sprintf("Post '%v' post does not exist.\n", filename))
    } else if exists == true {
        for _, tag := range tags.Tags() {
            AddTag(db, tag)
            fmt.Println("Add tag", tag)
        }
        query := fmt.Sprintf("INSERT INTO PostTag (post_id, post_tag) SELECT (SELECT id FROM Post WHERE filename = ?), id FROM Tag WHERE name IN (%v)", tags.String())
        if _, execErr := db.Exec(query, filename); execErr != nil {
            return execErr
        }
        log.Println("Query: ", query, filename)
    }
    return nil
}

func GetTagsFromPost(db *sql.DB, filename string) ([]string, error) {
    const QUERY = `SELECT name as tag FROM Tag
JOIN PostTag ON PostTag.post_tag = Tag.id
JOIN Post ON Post.id = PostTag.post_id
WHERE Post.filename = ?`
    rows, err := db.Query(QUERY, filename)
    defer rows.Close()
    if err != nil {
        return []string{}, err
    }
    tags := make([]string, 0)
    for rows.Next() == true {
        var tag string
        if err := rows.Scan(&tag); err != nil {
            return []string{}, err
        }
        tags = append(tags, tag)
    }
    return tags, nil
}

