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
    LIMIT %v OFFSET %v`;
    offset := pp.postsPerPage*pp.currentPage
    query := fmt.Sprintf(QUERY, pp.postsPerPage, offset)

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

    hasNext := pp.currentPage + 1 < pp.totalPages
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
}

func GetPostsPages(db *sql.DB, postsPerPage int64) (PostsPages, error) {
    const COUNT_QUERY = `SELECT COUNT(*) AS total_posts
FROM Post
WHERE id NOT IN
(SELECT PostTag.post_id FROM PostTag
JOIN Tag ON PostTag.tag_id = Tag.id
WHERE Tag.name = "page")`;

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
