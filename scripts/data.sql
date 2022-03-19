-- 建表语句
create table file
(
    id               TEXT not null
        constraint file_pk
            primary key,
    name             TEXT,
    path             TEXT,
    readme_url       TEXT,
    is_folder        integer default 0,
    download_url     TEXT,
    last_modify_time INTEGER,
    size             integer,
    password         TEXT,
    password_url     TEXT
);