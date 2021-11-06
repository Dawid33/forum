CREATE SCHEMA forum;
CREATE TABLE forum.users (
    userID      text,
    password    text
);
CREATE TABLE forum.posts (
    threadID    SERIAL UNIQUE PRIMARY KEY,
    userID      TEXT NOT NULL,
    title       TEXT,
    content     TEXT
);
CREATE TABLE forum.comments (
    commentID           SERIAL UNIQUE PRIMARY KEY,
    threadID            INT,
    parentID            INT,
    kidsID              INT[],
    userID              TEXT NOT NULL,
    content             TEXT,
    FOREIGN KEY (threadID) REFERENCES forum.posts(threadID)
);
