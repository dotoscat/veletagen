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
