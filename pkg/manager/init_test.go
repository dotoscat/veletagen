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
