package manager

import (
    "database/sql"
    "time"
    "fmt"
    "log"

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
    const QUERY = `SELECT id, filename, title, date FROM Post
    WHERE id NOT IN
    (SELECT PostTag.post_id FROM PostTag
    JOIN Tag ON PostTag.tag_id = Tag.id
    WHERE Tag.name = "page")
    LIMIT %v OFFSET %v`;
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
    id int64
    Filename string
    Title string
    Date time.Time
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
