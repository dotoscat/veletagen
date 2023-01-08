package constructor

import (
    // "time"
    "database/sql"
    "log"
    "embed"
    "path/filepath"
    "os"
    "text/template"
//    "io/fs"

    "github.com/dotoscat/veletagen/pkg/manager"
    "github.com/dotoscat/veletagen/pkg/common"
)

//go:embed templates/base.html templates/post.html
var postTemplate embed.FS

//go:embed templates/base.html templates/page.html
var pageTemplate embed.FS

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

func RenderTemplate(tmpl *template.Template, outputPath string) error {
        log.Println("Render template", tmpl);
        log.Println("to:", outputPath);
        outputFile, errFile := os.Create(outputPath)
        defer outputFile.Close()
        if errFile != nil {
            return errFile
        }
        if err := loadedBaseTemplate.Execute(outputFile, tmpl); err != nil {
            return err
        }
    return nil
}

func Construct(db *sql.DB, basePath string) error {
    var website common.WebsiteBase
    var err error
    if website, err = manager.GetWebsiteBase(db); err != nil {
        return err
    }
    log.Println("website base:", website)
    outputPath := website.OutputPath

    branches := []string{
            "posts",
            "pages",
    }
    common.CreateTree(outputPath, branches)

    indexPath := filepath.Join(outputPath, "index.html")

    var loadedPostTemplate *template.Template
    loadedPostTemplate, err = template.ParseFS(postTemplate, "templates/*")
    if err != nil {
        return err
    }

    if err := RenderTemplate(indexPath); err != nil {
        return err
    }

    return nil
}
