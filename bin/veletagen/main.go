package main

import (
    "database/sql"
    "flag"
    "log"
    "fmt"
    "strings"

    "github.com/dotoscat/veletagen/pkg/manager"
)

func main() {
    var init string
    var target string

    var getTitle bool
    var setTitle string

    var getPostsPerPage bool
    var setPostsPerPage int64

    var getLang bool
    var setLang string

    var addCSS string
    var removeCSS string
    var getCSS bool

    var addScript string
    var removeScript string
    var getScripts bool

    var addPost string
    var removePost string
    //var getPosts bool

    var post string
    var addTags Tags
    //var removeTags Tags

    flag.StringVar(&init, "init", "", "init <path>.")
    flag.StringVar(&target, "target", "", "target <path>.")

    flag.BoolVar(&getTitle, "get-title", false, "Get title used for building.")
    flag.StringVar(&setTitle, "set-title", "", "Set the title to be used for building.")

    flag.BoolVar(&getPostsPerPage, "get-posts-per-page", false, "Get the number of posts per page.")
    flag.Int64Var(&setPostsPerPage, "set-posts-per-page", 0, "Set the number of posts per page.")

    flag.BoolVar(&getLang, "get-lang", false, "Gets the site main lang.")
    flag.StringVar(&setLang, "set-lang", "", "Sets the site main lang.")

    flag.StringVar(&addCSS, "add-CSS", "", "Add a CSS file to be used for the whole website.")
    flag.StringVar(&removeCSS, "remove-CSS", "", "Remove a CSS file to be used for the whole website.")
    flag.BoolVar(&getCSS, "get-CSS", false, "Get CSS files added to the website.")

    flag.StringVar(&addScript, "add-script", "", "Add a JavaScript file to be used for the whole website.")
    flag.StringVar(&removeScript, "remove-script", "", "Remove a JavaScript file to be used for the whole website.")
    flag.BoolVar(&getScripts, "get-scripts", false, "Get JavaScript files added to the website.")

    flag.StringVar(&addPost, "add-post", "", "Add a post file to be used for the whole website.")
    flag.StringVar(&removePost, "remove-post", "", "Remove a post file to be used for the whole website.")

    flag.Var(&addTags, "add-tags", "Set an array of tags separated by ','")
    flag.StringVar(&post, "post", "", "Post to manipulate.")

    flag.Parse()

    if init != "" {
        log.Println("Init at: ", init)
        dbPath := manager.GetPathDB(init)
        log.Println("Index in", dbPath)
        if errCreateTree := manager.CreateTree(init); errCreateTree != nil {
            log.Fatal(errCreateTree)
        }
        if _, errOpenDatabase := manager.OpenDatabase(dbPath); errOpenDatabase != nil {
            log.Fatal(errOpenDatabase)
        }
        return
    }

    if target == "" {
        flag.PrintDefaults()
        return
    }

    var db *sql.DB
    var errOpenDatabase error

    db, errOpenDatabase = manager.OpenDatabase(manager.GetPathDB(target))
    defer db.Close()
    if errOpenDatabase != nil {
        log.Fatal(errOpenDatabase)
    }

    if setTitle != "" {
        if err := manager.SetTitle(db, setTitle); err != nil {
            log.Fatal(err)
        }
    }
    if getTitle == true {
        log.Println("Call function to get title to target: ", target)
        log.Println("Path DB: ", manager.GetPathDB(target))
        if title, err := manager.GetTitle(db); err != nil {
            log.Fatal(err)
        } else {
            fmt.Printf("title:%v\n", title)
        }
    }

    if setPostsPerPage > 0 {
        if err := manager.SetPostsPerPage(db, setPostsPerPage); err != nil {
            log.Fatal(err)
        }
    }
    if getPostsPerPage == true {
        if postsPerPage, err := manager.GetPostsPerPage(db); err != nil {
            log.Fatal(err)
        } else {
            fmt.Printf("posts-per-page:%v\n", postsPerPage)
        }
    }

    if setLang != "" {
        if err := manager.SetLang(db, setLang); err != nil {
            log.Fatal(err)
        }
    }
    if getLang == true {
        if lang, err := manager.GetLang(db); err != nil {
            log.Fatal(err)
        } else {
            fmt.Printf("lang:%v\n", lang)
        }
    }

    if addCSS != "" {
        if err := manager.AddCSS(db, addCSS); err != nil {
            log.Fatal(err)
        }
    } else if removeCSS != "" {
        if err := manager.RemoveCSS(db, removeCSS); err != nil {
            log.Fatal(err)
        }
    }

    if getCSS == true {
        if cssList, err := manager.GetCSS(db); err != nil {
            log.Fatal(err)
        } else {
            fmt.Printf("CSS:%v\n", strings.Join(cssList, ","))
        }
    }

    if addScript != "" {
        if err := manager.AddScript(db, addScript); err != nil {
            log.Fatal(err)
        }
    } else if removeScript != "" {
        if err := manager.RemoveScript(db, removeScript); err != nil {
            log.Fatal(err)
        }
    }

    if getScripts == true {
        if scriptList, err := manager.GetScripts(db); err != nil {
            log.Fatal(err)
        } else {
            fmt.Printf("Scripts:%v\n", strings.Join(scriptList, ","))
        }
    }

    if addPost != "" {
        if err := manager.AddPost(db, addPost); err != nil {
            log.Fatal(err)
        }
    } else if removePost != "" {
        if err := manager.RemovePost(db, removePost); err != nil {
            log.Fatal(err)
        }
    }

    if addTags.String() != "" {
        if post == "" {
            log.Println("Please, specify what post you want to add to.")
        } else {

        }
        log.Println(addTags)
    }

    log.Println("END")
}
