# Veletagen

This is a static website or blog generator. All the metadata is separated into a database that you manage it
from cli program. The content, markdown files, and assets are separated from the metadata in different folders.

## Init a project

    veletagen -init <path>

This creates the following tree at <path>:

/index.db
/posts
/assets/css
/assets/scripts

## Working with the project

Generally you will be working with the project this way:

    veletangen -target <path> <any other arguments>

### Posts

Post files, markdown files, must be put in /posts and must be
added to the database with

    veletagen -target <path> -add-post <filename>

So you can operate with them with categories and tags.

#### Categories

You can add and remove categories from posts.

    veletagen -target <path> -post <filename> -add-categories [category1, category2, ...]
    veletagen -target <path> -post <filename> -remove-categories [category1, category2, ...]
    veletagen -target <path> -post <filename> -get-categories

#### Tags

For the moment the only valid tag is "page" and is reserved to turn a post into a page.
You can add any tag to the posts but without any visible effect.

    veletagen -target <path> -post <filename> -add-tags [category1, category2, ...]
    veletagen -target <path> -post <filename> -remove-tags [category1, category2, ...]
    veletagen -target <path> -post <filename> -get-tags

### Assets

Assets must be placed inside /assets and then be managed by the database.

#### CSS

CSS must be stored at /assets/css

    veletagen -target <path> -add-css <file>
    veletagen -target <path> -remove-css <file>
    veletagen -target <path> -get-css <file>

#### Images

Images must be placed at /assets/css. There is not need to be managed to the database.
These files are copied at the output if used inside the post.

## Build project

This will construct the website only with the metadata added to the database.
It only counts the posts added to the database.

    veletagen -target <path> -build
