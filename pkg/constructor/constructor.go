package constructor

import (
    // "time"
    "database/sql"
    "log"
    "embed"
    "path/filepath"
    "os"
    "text/template"
    "io/fs"

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

func Construct(db *sql.DB, basePath string) error {
    var website common.WebsiteBase
    var err error
    if website, err = manager.GetWebsiteBase(db); err != nil {
        return err
    }
    log.Println("website base:", website)
    outputPath := website.OutputPath
    // var wd string
    // var wdErr error
    // wd, wdErr = os.Getwd()
    // if wdErr != nil {
    //    return wdErr
    // }

    // Create output structure
    // Base
    if err := os.MkdirAll(outputPath, fs.ModeDir); err != nil  {
        return err
    }
    // Posts
    postsPath := filepath.Join(outputPath, "posts")
    if err := os.MkdirAll(postsPath, fs.ModeDir); err != nil {
        return err
    }
    //

    indexPath := filepath.Join(outputPath, "index.html")
    if loadedBaseTemplate, err := template.ParseFS(baseTemplate, "templates/*"); err != nil {
        return err
    } else {
        log.Println("Create index file:", indexPath);
        indexFile, errFile := os.Create(indexPath)
        defer indexFile.Close()
        if errFile != nil {
            return errFile
        }
        if err := loadedBaseTemplate.Execute(indexFile, website); err != nil {
            return err
        }
    }

    return nil
}
