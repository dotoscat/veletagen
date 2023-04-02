package common

import (
    "path/filepath"
    "os"
    "io"
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

func CopyFile(src, dst string) error {
    var err error
    var srcFile *os.File
    var dstFile *os.File
    srcFile, err = os.Open(src)
    defer srcFile.Close()
    if err != nil {
        return err
    }
    dstFile, err = os.Create(dst)
    defer dstFile.Close()
    if err != nil {
        return err
    }
    if _ , err := io.Copy(dstFile, srcFile); err != nil {
        return err
    }
    return nil
}
