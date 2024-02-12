CREATE DATABASE feeds CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE feeds;

CREATE TABLE feed
(
    id              INTEGER             NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title           VARCHAR(300)        NOT NULL,
    name            VARCHAR(100) UNIQUE NOT NULL,
    url             VARCHAR(300) UNIQUE NOT NULL,
    description     TEXT,
    item_tag        VARCHAR(100)        NOT NULL,
    item_cls        VARCHAR(100)        NOT NULL,
    title_tag       VARCHAR(100)        NOT NULL,
    title_cls       VARCHAR(100),
    link_tag        VARCHAR(10)         NOT NULL,
    link_cls        VARCHAR(100),
    description_tag VARCHAR(100),
    description_cls VARCHAR(100),
    created         DATETIME            NOT NULL
);

CREATE INDEX idx_feed_name ON feed (name);

