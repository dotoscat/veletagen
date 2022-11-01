package main

import (
    "flag"
    "log"
    "path/filepath"

    "github.com/dotoscat/veletagen/pkg/manager"
)

func main() {
    var init string
    flag.StringVar(&init, "init", "", "init <path>")
    flag.Parse()
    if init != "" {
        log.Println("Init at: ", init)
        dbPath := filepath.Join(init, "index.db")
        log.Println("Index in", dbPath)
        if errCreateTree := manager.CreateTree(init); errCreateTree != nil {
            log.Fatal(errCreateTree)
        }
        if _, errOpenDatabase := manager.OpenDatabase(dbPath); errOpenDatabase != nil {
            log.Fatal(errOpenDatabase)
        }
    } else {
        flag.PrintDefaults()
    }
}
