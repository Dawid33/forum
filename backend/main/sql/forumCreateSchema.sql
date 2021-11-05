CREATE SCHEMA forum;
CREATE TABLE forum.users (
    userid      text,
    password    text
);
CREATE TABLE forum.posts (
    postid      SERIAL PRIMARY KEY,
    userid      TEXT NOT NULL ,
    post        TEXT
);