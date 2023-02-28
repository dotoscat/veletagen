package manager

import (
    "database/sql"
    "time"
    "fmt"

    // "github.com/dotoscat/veletagen/pkg/common"
)

func AddPost(db *sql.DB, filename string) error {
    return InsertStringInto(db, "Post", "filename", filename)
}

func RemovePost(db *sql.DB, filename string) error {
    return RemoveStringFrom(db, "Post", "filename", filename)
}

type PostsPages struct{
    currentPage int64
    totalPages int64
    postsPerPage int64
}

func (pp PostsPages) Next() bool {
    return pp.currentPage <= pp.totalPages
}

func (pp *PostsPages) GetPostsFromCurrentPage(db *sql.DB) (PostsPage, error) {
    const QUERY = `SELECT id, filename, title, date FROM Post LIMIT %v OFFSET %v`;
    offset := pp.postsPerPage*pp.currentPage
    query := fmt.Sprintf(QUERY, pp.postsPerPage, offset)
    fmt.Println(query)

    posts := make([]Post, 0)

    if rows, err := db.Query(query); err != nil {
        return PostsPage{}, err
    } else {
        defer rows.Close()
        for rows.Next() {
            var post Post
            var err error
            if post, err = CreatePostFromRows(rows); err != nil {
                return PostsPage{}, err
            }
            posts = append(posts, post)
        }
    }

    hasNext := pp.currentPage + 1 <= pp.totalPages
    hasPrevious := pp.currentPage - 1 >= 0

    postsPage := PostsPage{
        Number: pp.currentPage,
        Posts: posts,
        HasNext: hasNext,
        HasPrevious: hasPrevious,
    }
    pp.currentPage++
    return postsPage, nil
}

type Post struct {
    Name string
    Filename string
    Title string
    Date time.Time
}

func CreatePostFromRows(rows *sql.Rows) (Post, error) {
    var name string
    var filename string
    var title string
    var date time.Time

    if err := rows.Scan(&name, &filename, &title, &date); err != nil {
        return Post{}, err
    }
    return Post {name, filename, title, date}, nil
}

type PostsPage struct {
    Number int64
    Posts []Post
    HasNext bool
    HasPrevious bool
}

func GetPostsPages(db *sql.DB, postsPerPage int64) (PostsPages, error) {
    const COUNT_QUERY = `SELECT COUNT(*) AS total_posts
FROM Post
WHERE id NOT IN
(SELECT PostTag.post_id FROM PostTag
JOIN Tag ON PostTag.tag_id = Tag.id
WHERE Tag.name = "page")`;
    const QUERY = ``;

    postsPages := PostsPages{
        postsPerPage: postsPerPage,
    }
    var totalPosts int64
    row := db.QueryRow(COUNT_QUERY)
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

    return postsPages, nil
}
