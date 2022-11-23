package manager

import (
    "database/sql"
    //"errors"
    //"fmt"
)

func AddPost(db *sql.DB, filename string) error {
    return InsertStringInto(db, "Post", "filename", filename)
}

func RemovePost(db *sql.DB, filename string) error {
    return RemoveStringFrom(db, "Post", "filename", filename)
}

func AddTag(db *sql.DB, name string) error {
    if exists, err := DoesStringExistIn(db, "Tag", "name", name); err != nil {
        return err
    } else if exists == true {
        return nil
    }
    return InsertStringInto(db, "Tag", "name", name)
}

func RemoveTag(db *sql.DB, name string) error {
    return RemoveStringFrom(db, "Tag", "name", name)
}

/*
func AddTagsToPost(db *sql.DB, filename string, tags []string) error{
    if exists, err := DoesStringExistIn(db, "Post", "filename", filename); err != nil {
        return err
    } else exists == false {
        return errors.New(fmt.Sprintf("'%v' post does not exist.\n", filename))
    }
    query := fmt.Sprintf
    return true, nil
}
*/
