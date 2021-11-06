CREATE SCHEMA forum;
CREATE TABLE forum.users (
    userID      text,
    password    text
);
CREATE TABLE forum.posts (
    postID      SERIAL PRIMARY KEY,
    userID      TEXT NOT NULL,
    post        TEXT
);
CREATE TABLE forum.comments (
    postID              SERIAL PRIMARY KEY,
    userID              TEXT NOT NULL,
    parentCommentID     INT,
    commentID           SERIAL UNIQUE,
    comment             TEXT,
    childComments       INT[]
);
