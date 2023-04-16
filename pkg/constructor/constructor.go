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
//   limitations under the License.package constructor
package constructor

import (
    "database/sql"
    "log"
    "embed"
    "fmt"
    "strings"
    "path/filepath"
    "os"
    "text/template"
    "net/url"

    "github.com/dotoscat/veletagen/pkg/manager"
    "github.com/dotoscat/veletagen/pkg/common"

    "github.com/gomarkdown/markdown"
    "github.com/gomarkdown/markdown/ast"
)

//go:embed templates/base.html templates/post.html
var postTemplate embed.FS

//go:embed templates/base.html templates/postspage.html
var postsPageTemplate embed.FS

type Website struct {
    Config manager.Config
    basePath string
    Pages []PostWebpage
    Styles []string
    Categories map[string]string
    // scripts
}

func NewPostWebpageFromPost(post manager.Post, root string, website *Website) PostWebpage {
    filename, _ := strings.CutSuffix(post.Filename, ".md")
    postUrl := strings.Join([]string{root, filename + ".html"}, "/")
    webpage := NewWebpage(website, postUrl)
    srcPath := filepath.Join(website.basePath, root, post.Filename)

    imagesPath := make([]string, 0)

    if md, err := os.ReadFile(srcPath); err != nil {
        log.Println(err) // Maybe change this to return an error too?
    } else {
        root := markdown.Parse(md, nil)
        log.Println("Check whether", srcPath, "has images...")
        ast.WalkFunc(root, func(node ast.Node, entering bool) ast.WalkStatus {
            if entering == false {
                goto next
            }
            switch t := node.(type) {
                case *ast.Image:
                    src := string(t.Destination)
                    imagesPath = append(imagesPath, src)
                    break
            }
            next:
            return ast.GoToNext
        })

    }

    postWebpage := PostWebpage{
        Webpage: webpage,
        Post: post,
        src: srcPath,
        imagesPath: imagesPath,
    }
    return postWebpage
}

type Webpage struct {
    Website *Website
    Url string // Normalize to write output
    OutputPath string
}

func NewWebpage(website *Website, url string) Webpage {
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
    root string
}

func (ppw PostsPageWebpage) GetPreviousUrl() string {
    if ppw.PostsPage.HasPrevious == false {
       return ""
    }
    if ppw.PostsPage.Number - 1 <= 0  {
        return "/index.html"
    }
    path := strings.Join([]string{ppw.root, "page%v.html"}, "/")
    url := fmt.Sprintf(path, ppw.PostsPage.Number - 1 + 1)
    return url
}

func (ppw PostsPageWebpage) GetNextUrl() string {
    if ppw.PostsPage.HasNext == false {
        return ""
    }
    if ppw.PostsPage.Number == 0 { //index is page1
        return strings.Join([]string{ppw.root, "page2.html"}, "/")
    }
    path := strings.Join([]string{ppw.root, "page%v.html"}, "/")
    url := fmt.Sprintf(path, ppw.PostsPage.Number + 2)
    return url
}

func (ppw PostsPageWebpage) GetPreviousNumber() int64 {
    return ppw.PostsPage.Number - 1 + 1
}

func (ppw PostsPageWebpage) GetNextNumber() int64 {
    return ppw.PostsPage.Number + 1 + 1
}

func (ppw PostsPageWebpage) Number() int64 {
    return ppw.PostsPage.Number + 1
}

func NewPostsPageWebpage (website *Website, postsPage manager.PostsPage, root string) PostsPageWebpage {
    var url string
    if postsPage.Number == 0 {
        url = "index.html"
    } else {
        pageNumber := fmt.Sprintf("page%v.html", postsPage.Number + 1)
        url = strings.Join([]string{root, pageNumber}, "/")
    }
    webpage := NewWebpage(website, url)
    postWebpages := make([]PostWebpage, 0)
    for _, aPost := range postsPage.Posts {
        postWebpage := NewPostWebpageFromPost(aPost, "/posts", website)
        log.Println("postWebpage: ", postWebpage)
        postWebpages = append(postWebpages, postWebpage)
        // log.Println("aPost:", aPost)
    }
    // replace extension from filename for post output
    postsPageWebpage := PostsPageWebpage{
        webpage,
        postsPage,
        postWebpages,
        root,
    }
    //log.Println("webpage posts page:", webpage)
    return postsPageWebpage
}

type PostWebpage struct {
    Webpage
    Post manager.Post
    src string
    imagesPath []string
}

func (pw PostWebpage) ImagesPath() []string {
    return pw.imagesPath
}

func (pw PostWebpage) Content() string {
    // The idea here is to render markdown to html and give it as output
    if md, err := os.ReadFile(pw.src); err != nil {
        return err.Error()
    } else {
        html := markdown.ToHTML(md, nil, nil)
        return string(html)
    }
    return "Post content"
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

func BuildStylesPath(db *sql.DB, path string) ([]string, error) {
    paths := make([]string, 0)
    var styles []string
    var err error
    styles, err = manager.GetCSS(db)
    if err != nil {
        return paths, err
    }
    for _, style := range styles {
        path := strings.Join([]string{path, style}, "/")
        paths = append(paths, path)
    }
    return paths, nil
}

func Construct(db *sql.DB, basePath string) error {
    var config manager.Config
    var err error
    if config, err = manager.GetConfig(db); err != nil {
        return err
    }
    log.Println("config base:", config)

    var stylesPath []string

    stylesPath, err = BuildStylesPath(db, "/assets/css")
    if err != nil {
        return err
    }

    log.Println("stylesPath", stylesPath)

    website := Website{
        Config: config,
        basePath: basePath,
        Styles: stylesPath,
        Categories: make(map[string]string),
    }

    if pages, err := manager.GetPages(db); err != nil {
        return err
    } else {
        for _, page := range pages {
            pageWebpage := NewPostWebpageFromPost(page, "/posts", &website)
            website.Pages = append(website.Pages, pageWebpage)
        }
    }

    branches := []string{
            "posts",
            "pages",
            "assets/css",
            "assets/images",
    }

    var categories []string
    categories, err = manager.GetCategories(db)

    if err != nil {
        return err
    }
    for _, category := range categories {
        log.Println("category:", category)
        branches = append(branches, category)
        var categoryUrl string
        var err error
        categoryUrl, err = url.JoinPath("/", category, "index.html")
        if err != nil {
            return err
        }
        log.Println("category url:", categoryUrl)
        website.Categories[category] = categoryUrl
    }

    outputPath := config.OutputPath

    common.CreateTree(outputPath, branches)

    // Copy assets
    // styles
    for _, style := range website.Styles {
        srcStyle := filepath.Join(basePath, style)
        dstStyle := filepath.Join(outputPath, style)
        log.Println("Copy from:", srcStyle, ";to:", dstStyle)
        if err := common.CopyFile(srcStyle, dstStyle); err != nil {
            log.Fatal(err)
        }
    }

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

    for _, page := range website.Pages {
        if err := RenderTemplate(templates["post"], page.OutputPath, page); err != nil {
            return err
        }
    }
    postsPages, err := manager.GetPostsPages(db, config.PostsPerPage, "")
    for postsPages.Next() {
        if postsPage, err := postsPages.GetPostsFromCurrentPage(db); err != nil {
            return err
        } else {
            postsPageWebpage := NewPostsPageWebpage(&website, postsPage, "/pages")
            // log.Println("postsPageWebpage Number: ", postsPageWebpage.PostsPage.Number)
            // log.Println("postsPageWebpage HasPrevious: ", postsPageWebpage.PostsPage.HasPrevious)
            // log.Println("postsPageWebpage HasNext: ", postsPageWebpage.PostsPage.HasNext)
            // Render posts from postsPage
            log.Println("Posts from postsPageWebpage: ", postsPageWebpage)
            for _, post := range postsPageWebpage.Posts {
                if err := RenderTemplate(templates["post"], post.OutputPath, post); err != nil {
                    return err
                }
                // Copy post images if any
                for _, imagePath := range post.ImagesPath() {
                    log.Println("Copy this image:", imagePath)
                    imageSrc := filepath.Join(basePath, imagePath)
                    imageDst := filepath.Join(outputPath, imagePath)
                    if err := common.CopyFile(imageSrc, imageDst); err != nil {
                        log.Fatal(err)
                    }
                }
            }
            // Render postsPage
            if err := RenderTemplate(templates["postsPage"], postsPageWebpage.OutputPath, postsPageWebpage); err != nil {
                return err
            }
        }
    }

    // Render categories pages

    log.Println("website", website)
    log.Println("postsPerPage:", postsPages)

    return nil
}
