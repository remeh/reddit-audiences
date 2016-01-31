-- Database init
CREATE USER audiences WITH UNENCRYPTED PASSWORD 'audiences';
CREATE DATABASE "audiences";
GRANT ALL ON DATABASE "audiences" TO "audiences";

-- Switch to the audiences db as the audiences user.
\connect "audiences";
set role "audiences";

-- Tables
CREATE TABLE "subreddit" (
    "name" TEXT default '',
    "creation_time" TIMESTAMP WITH TIME ZONE,
    "last_crawl" TIMESTAMP WITH TIME ZONE,
    "next_crawl" TIMESTAMP WITH TIME ZONE,
    "active" BOOLEAN
);

CREATE INDEX ON "subreddit" ("next_crawl");

CREATE UNIQUE INDEX ON "subreddit" ("name");

CREATE TABLE "audience" (
    "subreddit" TEXT default '',
    "crawl_time" TIMESTAMP WITH TIME ZONE,
    "audience" INT
);

CREATE UNIQUE INDEX ON "audience" ("subreddit", "crawl_time");
CREATE INDEX ON "audience" ("subreddit", "crawl_time");
