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
    "time"
    "fmt"
    "log"

    // "github.com/dotoscat/veletagen/pkg/common"
)

func AddPost(db *sql.DB, filename, title string) error {
    const QUERY = "INSERT INTO Post (filename, title) VALUES (?, ?)"
    if _, err := db.Exec(QUERY, filename, title); err != nil {
        return err
    }
    return nil
}

func RemovePost(db *sql.DB, filename string) error {
    return RemoveStringFrom(db, "Post", "filename", filename)
}

type PostsPages struct{
    currentPage int64
    totalPages int64
    postsPerPage int64
    category string
}

func (pp PostsPages) Next() bool {
    return pp.currentPage < pp.totalPages
}

func (pp *PostsPages) GetPostsFromCurrentPage(db *sql.DB) (PostsPage, error) {
    const QUERY = `SELECT id, filename, title, date FROM Post
    WHERE id NOT IN
    (SELECT PostTag.post_id FROM PostTag
    JOIN Tag ON PostTag.tag_id = Tag.id
    WHERE Tag.name = "page")
    ORDER BY date DESC
    LIMIT %v OFFSET %v`
    const QUERY_CATEGORY = `SELECT Post.id, filename, title, date FROM Post
JOIN PostCategory
ON PostCategory.post_id = Post.id
JOIN Category
ON Category.id = PostCategory.category_id
WHERE Post.id NOT IN
(SELECT PostTag.post_id FROM PostTag
JOIN Tag ON PostTag.tag_id = Tag.id
WHERE Tag.name = "page")
AND Category.name = ?
ORDER BY date DESC
LIMIT %v OFFSET %v;`

    offset := pp.postsPerPage*pp.currentPage
    var query string
    var rows *sql.Rows
    var err error

    posts := make([]Post, 0)

    if pp.category == "" {
        query = fmt.Sprintf(QUERY, pp.postsPerPage, offset)
        rows, err = db.Query(query)
    } else {
        query = fmt.Sprintf(QUERY_CATEGORY, pp.postsPerPage, offset)
        rows, err = db.Query(query, pp.category)
    }
    defer rows.Close()

    if err != nil {
        return PostsPage{}, err
    }

    for rows.Next() {
        var post Post
        var err error
        if post, err = CreatePostFromRows(rows); err != nil {
            return PostsPage{}, err
        }
        posts = append(posts, post)
    }

    hasNext := pp.currentPage + 1 < pp.totalPages
    hasPrevious := pp.currentPage - 1 >= 0

    postsPage := PostsPage{
        Number: pp.currentPage,
        Posts: posts,
        HasNext: hasNext,
        HasPrevious: hasPrevious,
        Category: pp.category,
    }
    pp.currentPage++
    return postsPage, nil
}

type Post struct {
    id int64
    Filename string
    Title string
    Date time.Time
}

func GetPostByFilename (db *sql.DB, filename string) (Post, error) {
    const QUERY = `SELECT id, filename, title, date FROM Post WHERE filename = ?`
    post := Post{}
    row := db.QueryRow(QUERY, filename)
    if err := row.Err(); err != nil {
        return post, err
    }
    if err := row.Scan(&post.id, &post.Filename, &post.Title, &post.Date); err != nil {
        return post, err
    }
    return post, nil
}

func (p Post) Id() int64 {
    return p.id
}

func CreatePostFromRows(rows *sql.Rows) (Post, error) {
    var id int64
    var filename string
    var title string
    var date time.Time

    if err := rows.Scan(&id, &filename, &title, &date); err != nil {
        return Post{}, err
    }
    post := Post {id, filename, title, date}
    log.Println("CreatePostFromRows: ", post)

    return post, nil
}

func UpdatePostTitleByFilename(db *sql.DB, filename, title string) error {
    const QUERY = `UPDATE Post SET title = ? WHERE filename = ?`
    _, err := db.Exec(QUERY, title, filename)
    return err
}

type PostsPage struct {
    Number int64
    Posts []Post
    HasNext bool
    HasPrevious bool
    Category string
}

func GetPostsPages(db *sql.DB, postsPerPage int64, category string) (PostsPages, error) {
    const COUNT_QUERY = `SELECT COUNT(*) AS total_posts
FROM Post
WHERE id NOT IN
(SELECT PostTag.post_id FROM PostTag
JOIN Tag ON PostTag.tag_id = Tag.id
WHERE Tag.name = "page")`
    const COUNT_QUERY_CATEGORY = `
    SELECT COUNT(*) AS total_posts
FROM Post
JOIN PostCategory
ON PostCategory.post_id = Post.id
JOIN Category
ON Category.id = PostCategory.category_id
WHERE Post.id NOT IN
(SELECT PostTag.post_id FROM PostTag
JOIN Tag ON PostTag.tag_id = Tag.id
WHERE Tag.name = "page")
AND Category.name = ?;`

    query := COUNT_QUERY
    if category != "" {
        query = COUNT_QUERY_CATEGORY
    }

    postsPages := PostsPages{
        postsPerPage: postsPerPage,
        category: category,
    }
    var totalPosts int64
    var row *sql.Row

    if category == "" {
        row = db.QueryRow(query)
    } else {
        row = db.QueryRow(query, category)
    }

    if row.Err() != nil {
        return postsPages, row.Err()
    }
    if err := row.Scan(&totalPosts); err != nil {
        return postsPages, err
    }
    postsPages.totalPages = totalPosts / postsPerPage
    if totalPosts % postsPerPage > 0 {
        postsPages.totalPages++
    }
    // postsPages.postsPerPage = postsPerPage
    log.Println("totalPosts: ", totalPosts)
    log.Println("postsPages: ", postsPages)

    return postsPages, nil
}

// GetPages returns posts with the tag "page"
// These special posts are used as pages for navigation.
func GetPages(db *sql.DB) ([]Post, error) {
    const QUERY = `SELECT Post.id, filename, title, date FROM Post
JOIN PostTag ON PostTag.post_id = Post.id
JOIN Tag ON PostTag.tag_id = Tag.id
WHERE Tag.name = "page"`
    pages := make([]Post, 0)

    if rows, err := db.Query(QUERY); err != nil {
        return pages, err
    } else {
        defer rows.Close()
        for rows.Next() {
            if page, err := CreatePostFromRows(rows); err != nil {
                return pages, err
            } else {
                pages = append(pages, page)
            }
        }
    }

    return pages, nil
}
