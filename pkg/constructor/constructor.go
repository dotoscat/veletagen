package constructor

import (
    // "time"
    "database/sql"
    "log"
    "embed"
    "fmt"
    "strings"
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

type Website struct {
    Config manager.Config
    // categories, pages, scripts, styles...
}

type Webpage struct {
    Website Website
    Url string // Normalize to write output
    OutputPath string
}

func NewWebpage(website Website, url string) Webpage {
    return Webpage{
        Website: website,
        Url: url,
        OutputPath: filepath.Join(website.Config.OutputPath, url),
    }
}

type PostsPageWebpage struct {
    Webpage
    PostsPage manager.PostsPage
    Posts []PostWebpage
}

func NewPostsPageWebpage (website Website, postsPage manager.PostsPage) PostsPageWebpage {
    var url string
    if postsPage.Number == 0 {
        url = "index.html"
    } else {
        pageNumber := fmt.Sprintf("page%v.html", postsPage.Number + 1)
        url = strings.Join([]string{"/pages", pageNumber}, "/")
    }
    webpage := NewWebpage(website, url)
    postWebpages := make([]PostWebpage, 0)
    for _, aPost := range postsPage.Posts {
        //log.Println("aPost:", aPost)
        filename, _ := strings.CutSuffix(aPost.Filename, "md")
        postUrl := strings.Join([]string{"/posts", filename + "html"}, "/")
        webpage := NewWebpage(website, postUrl)
        postWebpage := PostWebpage{
            Webpage: webpage,
            Post: &aPost,
        }
        //log.Println("webPost:", postWebpage)
        postWebpages = append(postWebpages, postWebpage)
    }
    // replace extension from filename for post output
    postsPageWebpage := PostsPageWebpage{
        webpage,
        postsPage,
        postWebpages,
    }
    //log.Println("webpage posts page:", webpage)
    return postsPageWebpage
}

type PostWebpage struct {
    Webpage
    Post *manager.Post
}

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
    var config manager.Config
    var err error
    if config, err = manager.GetConfig(db); err != nil {
        return err
    }
    log.Println("config base:", config)

    website := Website{Config: config}

    outputPath := config.OutputPath

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

    if err := RenderTemplate(loadedPostTemplate, indexPath, config); err != nil {
        return err
    }

    postsPages, err := manager.GetPostsPages(db, 2) //TODO: Change that 2 by the one from the Config
    for postsPages.Next() {
        if postsPage, err := postsPages.GetPostsFromCurrentPage(db); err != nil {
            return err
        } else {
            postsPageWebpage := NewPostsPageWebpage(website, postsPage)
            log.Print("postsPageWebpage")
            log.Println(postsPageWebpage)
        }
    }
    log.Println("website", website)
    log.Println("postsPerPage:", postsPages)

    return nil
}
