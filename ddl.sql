drop table if exists feed;

create table feed
(
    id             integer primary key,
    title          text                                not null,
    name           text unique                         not null,
    url            text unique                         not null,
    description    text,
    item_selector  text                                not null,
    title_selector text                                not null,
    link_selector  text                                not null,
    desc_selector  text,
    created        TIMESTAMP DEFAULT CURRENT_TIMESTAMP not null
);

create index feed_name_idx on feed (name);
