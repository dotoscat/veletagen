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

func (pp PostsPages) HasNext() bool {
    return pp.currentPage < pp.totalPages;
}

func (pp PostsPages) HasLast() bool {
    return pp.currentPage > 0
}

func (pp *PostsPages) GoNext() bool {
    if pp.HasNext() {
        pp.currentPage++
        return true
    }
    return false
}

func (pp *PostsPages) GoLast() bool {
    if pp.HasLast() {
        pp.currentPage--
        return true
    }
    return false
}

func (pp PostsPages) CurrentPage() int64 {
    return pp.currentPage
}

func (pp PostsPages) GetPostsFromCurrentPage(db *sql.DB) PostPage {
    const QUERY = `SELECT id, filename, title, date FROM Post LIMIT %v OFFSET %v`;
    offset := pp.postsPerPage*pp.currentPage
    query := fmt.Sprintf(QUERY, pp.postsPerPage, offset)
    fmt.Println(query)
    return []Post{}
}

type Post struct {
    Name string
    Filename string
    Title string
    Date time.Time
}

type PostPage []Post

// func (pp PostsPages) GetPage(int64 page) PostPage {
//    return PostPage{}
//}

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
