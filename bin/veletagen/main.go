package main

import (
    "database/sql"
    "flag"
    "log"
    "fmt"

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

    flag.StringVar(&init, "init", "", "init <path>.")
    flag.StringVar(&target, "target", "", "target <path>.")
    flag.BoolVar(&getTitle, "get-title", false, "Get title used for building.")
    flag.StringVar(&setTitle, "set-title", "", "Set the title to be used for building.")
    flag.BoolVar(&getPostsPerPage, "get-posts-per-page", false, "Get the number of posts per page.")
    flag.Int64Var(&setPostsPerPage, "set-posts-per-page", 0, "Set the number of posts per page.")
    flag.BoolVar(&getLang, "get-lang", false, "Gets the site main lang.")
    flag.StringVar(&setLang, "set-lang", "", "Sets the site main lang.")

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

    log.Println("END")
}
