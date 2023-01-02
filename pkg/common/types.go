package common

type Category struct{
    Name string
    Url string
}

type WebsiteBase struct {
    Title string
    Categories []Category
    License string
    Lang string
    OutputPath string
    PostsPerPage int64
    // styles, scripts, ...
}
