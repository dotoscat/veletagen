package constructor

import (
    // "time"
    "database/sql"
    "log"
    "embed"

    "github.com/dotoscat/veletagen/pkg/manager"
    "github.com/dotoscat/veletagen/pkg/common"
)

//go:embed templates/base.html
var baseTemplate embed.FS

/*
type Webpage struct {
    Base WebsiteBase
    Output string
    Url string
}

type Post struct {
    Name string
    Filename string
    Title string
    Date time.Time
}

type PostsPage struct {
    Webpage
    LastPage *PostsPage
    NextPage *PostsPage
    Posts []Post
}
*/

func Construct(db *sql.DB) error {
    var website common.WebsiteBase
    var err error
    if website, err = manager.GetWebsiteBase(db); err != nil {
        return err
    }
    log.Println("website base:", website)

    return err
}
