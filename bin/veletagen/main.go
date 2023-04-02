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
//   limitations under the License.

package main

import (
    "database/sql"
    "flag"
    "log"
    "fmt"
    "strings"

    "github.com/dotoscat/veletagen/pkg/common"
    "github.com/dotoscat/veletagen/pkg/manager"
    "github.com/dotoscat/veletagen/pkg/constructor"
)

func main() {
    var build bool

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

    var addTags common.Tags
    var getTags bool
    var removeTags common.Tags

    var addCategories common.Tags
    var getCategories bool
    var removeCategories common.Tags

    flag.BoolVar(&build, "build", false, "Start building the site specified by target.")

    flag.StringVar(&init, "init", "", "init <path>.")
    flag.StringVar(&target, "target", "", "target <path>.")

    flag.BoolVar(&getTitle, "get-title", false, "Get title used for building or post.")
    flag.StringVar(&setTitle, "set-title", "", "Set the title to be used for building or post.")

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

    flag.StringVar(&post, "post", "", "Post to manipulate.")

    flag.Var(&addTags, "add-tags", "Set an array of tags separated by ','")
    flag.Var(&removeTags, "remove-tags", "Remove an array of tags separated by ','")
    flag.BoolVar(&getTags, "get-tags", false, "Gets tags related with this post.")

    flag.Var(&addCategories, "add-categories", "Add categories separated by ','")
    flag.Var(&removeCategories, "remove-categories", "Remove categories separated by ','")
    flag.BoolVar(&getCategories, "get-categories", false, "Gets categories related with this post.")

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
        log.Println("Don't forget the target. Use -target=<path>")
        return
    }

    var db *sql.DB
    var errOpenDatabase error

    db, errOpenDatabase = manager.OpenDatabase(manager.GetPathDB(target))
    defer db.Close()

    if errOpenDatabase != nil {
        log.Fatal(errOpenDatabase)
    }

    if build == true {
        log.Println("Start building the site!");
        if err := constructor.Construct(db, target); err != nil {
            log.Fatal(err)
        }
        return
    }

    if setTitle != "" && post == "" {
        if err := manager.SetTitle(db, setTitle); err != nil {
            log.Fatal(err)
        }
    } else if setTitle != "" && post != "" {
        if err := manager.UpdatePostTitleByFilename(db, post, setTitle); err != nil {
            log.Fatal(err)
        }
    }

    if getTitle == true && post == "" {
        log.Println("Call function to get title to target: ", target)
        log.Println("Path DB: ", manager.GetPathDB(target))
        if title, err := manager.GetTitle(db); err != nil {
            log.Fatal(err)
        } else {
            fmt.Printf("title:%v\n", title)
        }
    } else if getTitle == true && post != "" {
        if postObject, err := manager.GetPostByFilename(db, post); err != nil {
            log.Fatal(err)
        } else {
            fmt.Printf("title:%v\n", postObject.Title)
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
        if setTitle == "" {
            fmt.Printf("Please, set a title for post '%v'", addPost)
            return
        }
        if err := manager.AddPost(db, addPost, setTitle); err != nil {
            log.Fatal(err)
        }
    } else if removePost != "" {
        if err := manager.RemovePost(db, removePost); err != nil {
            log.Fatal(err)
        }
    }

    if addTags.String() != "" {
        if post == "" {
            log.Fatal("Please, specify what post you want to add to.")
        } else if err := manager.AddTagsToPost(db, post, addTags); err != nil {
            log.Fatal(err)
        }
        log.Println(addTags.String())
    } else if removeTags.String() != "" {
        if post == "" {
            log.Fatal("Please, specify what post you want to remove tags from.")
        } else if err := manager.RemoveTagsFromPost(db, post, removeTags); err != nil {
            log.Fatal(err)
        }
    }
    if getTags == true {
        if post == "" {
            log.Fatal("Please, specify what post you want to get tags from.")
        } else if tags, err := manager.GetTagsFromPost(db, post); err != nil{
            log.Fatal(err)
        } else {
            fmt.Println(tags)
        }
    }

    if addCategories.String() != "" {
        if post == "" {
            log.Fatal("Please, specify what post you want to add to categories.")
        } else if err := manager.AddCategoriesToPost(db, post, addCategories); err != nil {
            log.Fatal(err)
        }
    } else if removeCategories.String() != "" {
        if post == "" {
            log.Fatal("Please, specify what post you want to remove categories from.")
        } else if err := manager.RemoveCategoriesFromPost(db, post, removeCategories); err != nil {
            log.Fatal(err)
        }
    }
    if getCategories == true {
        if post == "" {
            log.Fatal("Plase, specify what post you want to get categories from.")
        } else if categories, err := manager.GetCategoriesFromPost(db, post); err != nil {
            log.Fatal(err)
        } else {
            fmt.Println(categories)
        }
    }

    log.Println("END")
}
