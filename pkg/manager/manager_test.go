package manager

import (
    "testing"
    "os"
    "flag"
    "log"
)

const GOLDEN_DATABASE = "./testdata/golden-database.db"

func TestMain (m *testing.M) {
    var goldenDatabase bool
    flag.BoolVar(&goldenDatabase, "generate-golden-database", false, "Re-generate an empty database to be used as golden file for the testing.")
    flag.Parse()
    if goldenDatabase == true {
        os.Remove(GOLDEN_DATABASE)
        if _, err := OpenDatabase(GOLDEN_DATABASE); err != nil {
            log.Fatal(err)
        }
        os.Exit(0)
    }
    os.Exit(m.Run())
}
