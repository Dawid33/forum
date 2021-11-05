CREATE SCHEMA forum;
CREATE TABLE forum.users (
     userid      text,
     password    text
);
CREATE TABLE forum.posts (
     userid      text,
     postid      text primary key,
     postDate    date,
     post        text
);