create database test;

use test;

create table tag_tbl(
    id int not null auto_increment,
    name varchar(40) not null,
    created_at datetime not null default current_timestamp,
    primary key (`id`),
    unique key `name` (`name`) using hash
)engine=InnoDB default charset=utf8mb4;

create table entity_tag_tbl(
    id int(10) unsigned not null auto_increment,
    entity_id int(10) unsigned not null,
    tag_id int(10) unsigned not null,
    created_at datetime not null default current_timestamp,
    primary key (`id`),
    unique key `entity_id` (`entity_id`,`tag_id`) using btree
)engine=InnoDB default charset=utf8mb4;