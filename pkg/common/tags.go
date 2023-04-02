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
//   limitations under the License.package common
package common

import (
    "strings"
    "fmt"
)

type Tags struct {
    tags []string
}

func (t *Tags) String() string {
    wrappedTags := make([]string, len(t.tags))
    for i, tag := range t.tags {
        wrappedTags[i] = fmt.Sprintf("'%v'",tag)
    }
    return strings.Join(wrappedTags, ",")
}

func (t *Tags) Set (value string) error {
    elements := strings.Split(value, ",")
    tags := make([]string, len(elements))
    for i, element := range elements {
        tags[i] = strings.TrimSpace(element)
    }
    t.tags = tags
    return nil
}

func (t *Tags) Tags() []string {
    return t.tags
}
