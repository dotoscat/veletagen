package manager

import (
    "database/sql"
    "errors"
    "fmt"
    "log"

    "github.com/dotoscat/veletagen/pkg/common"
)

func AddCategory(db *sql.DB, name string) error {
    if exists, err := DoesStringExistIn(db, "Category", "name", name); err != nil {
        return err
    } else if exists == true {
        return nil
    }
    return InsertStringInto(db, "Category", "name", name)
}

func RemoveCategory(db *sql.DB, name string) error {
    return RemoveStringFrom(db, "Category", "name", name)
}

func AddCategoriesToPost(db *sql.DB, filename string, categories common.Tags) error{
    if exists, err := DoesStringExistIn(db, "Post", "filename", filename); err != nil {
        return err
    } else if exists == false {
        return errors.New(fmt.Sprintf("Post '%v' post does not exist.\n", filename))
    } else if exists == true {
        for _, category := range categories.Tags() {
            AddCategory(db, category)
            fmt.Println("Add category", category)
        }
        query := fmt.Sprintf("INSERT INTO PostCategory (post_id, category_id) SELECT (SELECT id FROM Post WHERE filename = ?), id FROM Category WHERE name IN (%v)", categories.String())
        if _, execErr := db.Exec(query, filename); execErr != nil {
            return execErr
        }
        log.Println("Query: ", query, filename)
    }
    return nil
}

func RemoveCategoriesFromPost(db *sql.DB, filename string, tags common.Tags) error {
    if exists, err := DoesStringExistIn(db, "Post", "filename", filename); err != nil {
        return err
    } else if exists == false {
        return errors.New(fmt.Sprintf("Post '%v' post does not exist.\n", filename))
    } else if exists == true {
        query := fmt.Sprintf(`DELETE FROM PostCategory WHERE PostCategory.id IN
            (SELECT PostCategory.id FROM PostCategory
            JOIN Post ON Post.id = PostCategory.post_id
            JOIN Category ON Category.id = PostCategory.category_id
            WHERE Post.filename = ? AND
            Category.name IN (%v));`, tags.String())
        if _, execErr := db.Exec(query, filename); execErr != nil {
            return execErr
        }
    }
    return nil
}

func GetCategoriesFromPost(db *sql.DB, filename string) ([]string, error) {
    const QUERY = `SELECT name as category FROM Category
JOIN PostCategory ON PostCategory.category_id = Category.id
JOIN Post ON Post.id = PostCategory.post_id
WHERE Post.filename = ?`
    rows, err := db.Query(QUERY, filename)
    defer rows.Close()
    if err != nil {
        return []string{}, err
    }
    categories := make([]string, 0)
    for rows.Next() == true {
        var category string
        if err := rows.Scan(&category); err != nil {
            return []string{}, err
        }
        categories = append(categories, category)
    }
    return categories, nil
}
