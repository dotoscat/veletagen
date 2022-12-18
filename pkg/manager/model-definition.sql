--
-- File generated with SQLiteStudio v3.3.3 on sá. sep. 10 12:37:33 2022
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: Category
CREATE TABLE IF NOT EXISTS Category (
    id   INTEGER PRIMARY KEY,
    name TEXT    UNIQUE
                 NOT NULL
);

-- Table: Config
CREATE TABLE IF NOT EXISTS Config (
    version        INTEGER DEFAULT (0)
                           NOT NULL,
    title          TEXT    DEFAULT [Your Website]
                           NOT NULL,
    posts_per_page INTEGER DEFAULT (3)
                           NOT NULL,
    output_path    TEXT    DEFAULT output
                           NOT NULL,
    lang           TEXT    DEFAULT en
                           NOT NULL
);

INSERT INTO Config (
                       version,
                       title,
                       posts_per_page,
                       output_path,
                       lang
                   )
                   VALUES (
                       0,
                       'Your Website',
                       3,
                       'output',
                       'en'
                   );

-- Table: ConfigCSS
CREATE TABLE IF NOT EXISTS ConfigCSS (
    id       INTEGER PRIMARY KEY,
    filename TEXT    UNIQUE
                     NOT NULL
);


-- Table: ConfigScript
CREATE TABLE IF NOT EXISTS ConfigScript (
    id       INTEGER PRIMARY KEY,
    filename TEXT    NOT NULL
                     UNIQUE
);


-- Table: Post
CREATE TABLE IF NOT EXISTS Post (
    id       INTEGER  PRIMARY KEY,
    filename TEXT     UNIQUE
                      NOT NULL,
    title    TEXT       DEFAULT [Post title]
                        NOT NULL,
    date     DATETIME NOT NULL
                      DEFAULT (CURRENT_DATE) 
);


-- Table: PostCategory
CREATE TABLE IF NOT EXISTS PostCategory (
    post_id       INTEGER REFERENCES Post (id) ON DELETE CASCADE
                                               ON UPDATE CASCADE,
    post_category INTEGER REFERENCES Category (id) ON DELETE CASCADE
                                                   ON UPDATE CASCADE
);


-- Table: PostCSS
CREATE TABLE IF NOT EXISTS PostCSS (
    filename TEXT    UNIQUE
                     NOT NULL,
    post_id  INTEGER REFERENCES Post (id) ON DELETE CASCADE
                                          ON UPDATE CASCADE
);


-- Table: PostOverridingCSS
CREATE TABLE IF NOT EXISTS PostOverridingCSS (
    filename TEXT    UNIQUE
                     NOT NULL,
    post_id  INTEGER REFERENCES Post (id) ON DELETE CASCADE
                                          ON UPDATE CASCADE
);


-- Table: PostOverridingScript
CREATE TABLE IF NOT EXISTS PostOverridingScript (
    filename TEXT    UNIQUE
                     NOT NULL,
    post_id  INTEGER REFERENCES Post (id) ON DELETE CASCADE
                                          ON UPDATE CASCADE
);


-- Table: PostScript
CREATE TABLE IF NOT EXISTS PostScript (
    filename TEXT    UNIQUE
                     NOT NULL,
    post_id  INTEGER REFERENCES Post (id) ON DELETE CASCADE
                                          ON UPDATE CASCADE
);


-- Table: PostTag
CREATE TABLE IF NOT EXISTS PostTag (
    id       INTEGER PRIMARY KEY,
    post_id  INTEGER REFERENCES Post (id) ON DELETE CASCADE
                                          ON UPDATE CASCADE,
    post_tag INTEGER REFERENCES Tag (id) ON DELETE CASCADE
                                         ON UPDATE CASCADE,
    UNIQUE(post_id, post_tag)
);


-- Table: Tag
CREATE TABLE IF NOT EXISTS Tag (
    id   INTEGER PRIMARY KEY,
    name TEXT    UNIQUE
                 NOT NULL
);


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
