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
//   limitations under the License.package manager
package manager

import (
    "testing"
    "os"
    "path/filepath"
    "crypto/sha256"
    "log"
)

func sameCheckSum(aPath, bPath string) (bool, error) {
    aBuffer, aErr := os.ReadFile(aPath)
    if aErr != nil {
        return false, aErr
    }
    bBuffer, bErr := os.ReadFile(bPath)
    if bErr != nil {
        return false, bErr
    }
    aSum := sha256.Sum256(aBuffer)
    bSum := sha256.Sum256(bBuffer)
    log.Print(aSum)
    log.Print(bSum)
    return aSum == bSum, nil
}

func TestOpenDatabase(t *testing.T) {
    tmpDir := os.TempDir()
    tempDatabase := filepath.Join(tmpDir, "test.db")
    t.Log(tempDatabase)
    os.Remove(tempDatabase)
    if db, err := OpenDatabase(tempDatabase); err != nil {
        t.Fatal(err)
    } else {
        t.Log(db)
        db.Close()
    }
    if same, err := sameCheckSum(tempDatabase, GOLDEN_DATABASE); err != nil {
        t.Fatal(err)
    } else if same == false {
        t.Fatal("They are not the same database.")
    }
    os.Remove(tempDatabase)
}

func TestCreateTree(t *testing.T) {
    tempBaseDir, err := os.MkdirTemp("", "test*")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(tempBaseDir)
    t.Log(tempBaseDir)
    CreateTree(tempBaseDir)
    postsPath := filepath.Join(tempBaseDir, "posts")
    scriptsPath := filepath.Join(tempBaseDir, SCRIPTS_PATH)
    cssPath := filepath.Join(tempBaseDir, CSS_PATH)
    if info, err := os.Lstat(postsPath); err != nil {
        t.Fatal(err)
    } else if info.IsDir() == false {
        t.Fatal(postsPath + " is not a dir")
    }
    if info, err := os.Lstat(scriptsPath); err != nil {
        t.Fatal(err)
    } else if info.IsDir() == false {
        t.Fatal(scriptsPath + " is not a dir")
    }
    if info, err := os.Lstat(cssPath); err != nil {
        t.Fatal(err)
    } else if info.IsDir() == false {
        t.Fatal(cssPath + " is not a dir")
    }
}
