-- Database init
CREATE USER audiences WITH UNENCRYPTED PASSWORD 'audiences';
CREATE DATABASE "audiences";
GRANT ALL ON DATABASE "audiences" TO "audiences";

-- Switch to the audiences db as the audiences user.
\connect "audiences";
set role "audiences";

-- Tables

-- subreddit

CREATE TABLE "subreddit" (
    "name" TEXT default '',
    "creation_time" TIMESTAMP WITH TIME ZONE,
    "last_crawl" TIMESTAMP WITH TIME ZONE,
    "active" BOOLEAN
);

CREATE INDEX ON "subreddit" ("last_crawl");

CREATE UNIQUE INDEX ON "subreddit" ("name");

-- audience

CREATE TABLE "audience" (
    "subreddit" TEXT default '',
    "crawl_time" TIMESTAMP WITH TIME ZONE,
    "audience" INT DEFAULT 0,
    "subscribers" INT DEFAULT 0
);

CREATE UNIQUE INDEX ON "audience" ("subreddit", "crawl_time");
CREATE INDEX ON "audience" ("subreddit", "crawl_time");

-- article

CREATE TABLE "article" (
    "subreddit" TEXT DEFAULT '', -- foreign key to subreddit
    "article_id" TEXT DEFAULT '', -- reddit article id
    "article_title" TEXT DEFAULT '',
    "article_external_link" TEXT DEFAULT '',
    "article_link" TEXT DEFAULT '',
    "author" TEXT DEFAULT '',
    "rank" INT DEFAULT 0,
    "crawl_time" TIMESTAMP WITH TIME ZONE,
    "promoted" BOOLEAN DEFAULT false,
    "sticky" BOOLEAN DEFAULT false
);

CREATE INDEX ON "article" ("subreddit", "crawl_time");

-- user

CREATE TABLE "user" (
    "uuid" TEXT DEFAULT '',
    "email" TEXT DEFAULT '',
    "hash" TEXT DEFAULT '',
    "firstname" TEXT default '',
    "lastname" TEXT default '',
    "creation_time" TIMESTAMP WITH TIME ZONE,
    "last_login" TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ON "user" ("uuid");
CREATE UNIQUE INDEX ON "user" ("email");

-- session

CREATE TABLE "session" (
    "token" TEXT DEFAULT '',
    "uuid" TEXT DEFAULT '',
    "hit_time" TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX ON "session" ("token");
ALTER TABLE "session" ADD FOREIGN KEY ("uuid") REFERENCES "user" ("uuid");

-- annotation

CREATE TABLE "annotation" (
    "owner" TEXT NOT NULL,
    "subreddit" TEXT NOT NULL,
    "time" TIMESTAMP WITH TIME ZONE, 
    "message" TEXT default ''
);

CREATE INDEX ON "annotation" ("owner", "subreddit", "time");
ALTER TABLE "annotation" ADD FOREIGN KEY ("owner") REFERENCES "user" ("uuid");
