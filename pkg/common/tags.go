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
