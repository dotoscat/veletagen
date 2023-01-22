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
var postsPageTemplate embed.FS

/*
type Webpage struct {
    Base WebsiteBase
    Output string
    Url string
}



type PostsPage struct {
    Webpage
    LastPage *PostsPage
    NextPage *PostsPage
    Posts []Post
}
*/

func RenderTemplate(tmpl *template.Template, outputPath string, data any) error {
        log.Println("Render template", tmpl);
        log.Println("to:", outputPath);
        outputFile, errFile := os.Create(outputPath)
        defer outputFile.Close()
        if errFile != nil {
            return errFile
        }
        if err := tmpl.Execute(outputFile, data); err != nil {
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

    // Load templates
    templates := make(map[string]*template.Template)
    templatesDefinition := []struct{
        name string
        fs embed.FS
    }{
        {"post", postTemplate},
        {"postsPage", postsPageTemplate},
    }

    for _, tuple := range templatesDefinition {
        if loadedTemplate, err := template.ParseFS(tuple.fs, "templates/*"); err != nil {
            return err
        } else {
            templates[tuple.name] = loadedTemplate
        }
        log.Println("Load template:", tuple)
    }
    log.Println("templates", templates)
    // End loading templates

    indexPath := filepath.Join(outputPath, "index.html")

    var loadedPostTemplate *template.Template
    loadedPostTemplate, err = template.ParseFS(postTemplate, "templates/*")
    if err != nil {
        return err
    }

    if err := RenderTemplate(loadedPostTemplate, indexPath, website); err != nil {
        return err
    }

    postsPerPage, err := manager.GetPostsPages(db, 2)
    for postsPerPage.Next() {
        if postsPages, err := postsPerPage.GetPostsFromCurrentPage(db); err != nil {
            return err
        } else {
            log.Println(postsPages)
        }
    }
    log.Println("postsPerPage:", postsPerPage)

    return nil
}
