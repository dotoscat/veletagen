package common

import (
    "path/filepath"
    "os"
    "io/fs"
    "log"
)

func CreateTree(base string, branches []string) error {
    for _, branch := range branches {
        path := filepath.Join(base, branch)
        log.Println("branch:", path)
        if err := os.MkdirAll(path, fs.ModeDir); err != nil {
            return err;
        }
    }
    return nil
}
